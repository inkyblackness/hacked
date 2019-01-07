package messages

import (
	"fmt"

	"github.com/inkyblackness/imgui-go"

	"github.com/inkyblackness/hacked/editor/external"
	"github.com/inkyblackness/hacked/editor/graphics"
	"github.com/inkyblackness/hacked/editor/model"
	"github.com/inkyblackness/hacked/editor/render"
	"github.com/inkyblackness/hacked/editor/values"
	"github.com/inkyblackness/hacked/ss1/content/audio"
	"github.com/inkyblackness/hacked/ss1/content/movie"
	"github.com/inkyblackness/hacked/ss1/content/text"
	"github.com/inkyblackness/hacked/ss1/edit/undoable/cmd"
	"github.com/inkyblackness/hacked/ss1/resource"
	"github.com/inkyblackness/hacked/ss1/world/ids"
	"github.com/inkyblackness/hacked/ui/gui"
)

type messageInfo struct {
	title string
}

var knownMessageTypes = map[resource.ID]messageInfo{
	ids.MailsStart:     {title: "Mails"},
	ids.LogsStart:      {title: "Logs"},
	ids.FragmentsStart: {title: "Fragments"},
}
var knownMessageTypesOrder = []resource.ID{ids.MailsStart, ids.LogsStart, ids.FragmentsStart}

var knownMailTypes = []string{
	"Audio/Text",
	"VMail: 'Shield' (0)",
	"VMail: 'Grove' (1)",
	"VMail: 'Bridge' (2)",
	"VMail: 'Laser' (3, 4)",
	"VMail: 'Status' (11)",
	"VMail: 'Explode' (5..9)",
}

const videoMailDisplayBase = 256

// View provides edit controls for messages.
type View struct {
	mod          *model.Mod
	messageCache *text.ElectronicMessageCache
	cp           text.Codepage
	movieCache   *movie.Cache
	imageCache   *graphics.TextureCache

	modalStateMachine gui.ModalStateMachine
	clipboard         external.Clipboard
	guiScale          float32
	commander         cmd.Commander

	model viewModel
}

// NewMessagesView returns a new instance.
func NewMessagesView(mod *model.Mod, messageCache *text.ElectronicMessageCache, cp text.Codepage,
	movieCache *movie.Cache, imageCache *graphics.TextureCache,
	modalStateMachine gui.ModalStateMachine, clipboard external.Clipboard,
	guiScale float32, commander cmd.Commander) *View {
	view := &View{
		mod:          mod,
		messageCache: messageCache,
		cp:           cp,
		movieCache:   movieCache,
		imageCache:   imageCache,

		modalStateMachine: modalStateMachine,
		clipboard:         clipboard,
		guiScale:          guiScale,
		commander:         commander,

		model: freshViewModel(),
	}
	return view
}

// WindowOpen returns the flag address, to be used with the main menu.
func (view *View) WindowOpen() *bool {
	return &view.model.windowOpen
}

// Render renders the view.
func (view *View) Render() {
	if view.model.restoreFocus {
		imgui.SetNextWindowFocus()
		view.model.restoreFocus = false
		view.model.windowOpen = true
	}
	if view.model.windowOpen {
		imgui.SetNextWindowSizeV(imgui.Vec2{X: 640 * view.guiScale, Y: 480 * view.guiScale}, imgui.ConditionOnce)
		if imgui.BeginV("Messages", view.WindowOpen(), imgui.WindowFlagsNoCollapse) {
			view.renderContent()
		}
		imgui.End()
	}
}

func (view *View) renderContent() {
	if imgui.BeginChildV("Properties", imgui.Vec2{X: 350 * view.guiScale, Y: -200 * view.guiScale}, false, 0) {
		imgui.PushItemWidth(-150 * view.guiScale)
		if imgui.BeginCombo("Message Type", knownMessageTypes[view.model.currentKey.ID].title) {
			for _, id := range knownMessageTypesOrder {
				if imgui.SelectableV(knownMessageTypes[id].title, id == view.model.currentKey.ID, 0, imgui.Vec2{}) {
					view.model.currentKey.ID = id
					view.model.currentKey.Index = 0
				}
			}
			imgui.EndCombo()
		}
		info, _ := ids.Info(view.model.currentKey.ID)
		index := view.model.currentKey.Index
		if gui.StepSliderInt("Index", &index, 0, info.MaxCount-1) {
			view.model.currentKey.Index = index
		}
		if imgui.BeginCombo("Language", view.model.currentKey.Lang.String()) {
			languages := resource.Languages()
			for _, lang := range languages {
				if imgui.SelectableV(lang.String(), lang == view.model.currentKey.Lang, 0, imgui.Vec2{}) {
					view.model.currentKey.Lang = lang
				}
			}
			imgui.EndCombo()
		}

		message, readOnly := view.currentMessage()
		availableDisplays := view.availableDisplays()

		if imgui.ButtonV("Clear", imgui.Vec2{}) {
			view.requestClear()
		}
		imgui.SameLine()
		if !readOnly {
			if imgui.ButtonV("Remove", imgui.Vec2{}) {
				view.requestRemove()
			}
		} else {
			imgui.Text("(read-only)")
		}
		if view.hasAudio() {
			imgui.Separator()
			sound := view.currentSound()
			if !sound.Empty() {
				imgui.LabelText("Audio", fmt.Sprintf("%.2f sec", sound.Duration()))
				if imgui.Button("Export") {
					view.requestExportAudio(sound)
				}
				imgui.SameLine()
			} else {
				imgui.LabelText("Audio", "(no sound)")
			}
			if imgui.Button("Import") {
				view.requestImportAudio()
			}
		}
		imgui.Separator()

		textModes := map[bool]string{true: "Verbose", false: "Terse"}
		if imgui.BeginCombo("Text Mode", textModes[view.model.showVerboseText]) {
			for _, setting := range []bool{true, false} {
				if imgui.SelectableV(textModes[setting], view.model.showVerboseText == setting, 0, imgui.Vec2{}) {
					view.model.showVerboseText = setting
				}
			}
			imgui.EndCombo()
		}
		view.renderText(readOnly, "Title", message.Title, func(newValue string) {
			view.requestTextChange(func(msg *text.ElectronicMessage) { msg.Title = text.Blocked(newValue)[0] })
		})
		view.renderText(readOnly, "Sender", message.Sender, func(newValue string) {
			view.requestTextChange(func(msg *text.ElectronicMessage) { msg.Sender = text.Blocked(newValue)[0] })
		})
		view.renderText(readOnly, "Subject", message.Subject, func(newValue string) {
			view.requestTextChange(func(msg *text.ElectronicMessage) { msg.Subject = text.Blocked(newValue)[0] })
		})
		nextMessageIndex := values.NewUnifier()
		nextMessageIndex.Add(message.NextMessage)
		values.RenderUnifiedSliderInt(readOnly, false, "Next Message Index", nextMessageIndex,
			func(u values.Unifier) int { return u.Unified().(int) },
			func(int) string { return "%d" },
			-1, 255,
			func(newValue int) {
				view.requestPropertyChange(func(msg *text.ElectronicMessage) { msg.NextMessage = newValue })
			})

		isInterrupt := values.NewUnifier()
		isInterrupt.Add(message.IsInterrupt)
		values.RenderUnifiedCheckboxCombo(readOnly, false, "Is Interrupt", isInterrupt,
			func(newValue bool) {
				view.requestPropertyChange(func(msg *text.ElectronicMessage) { msg.IsInterrupt = newValue })
			})

		colorIndex := values.NewUnifier()
		colorIndex.Add(message.ColorIndex)
		values.RenderUnifiedSliderInt(readOnly, false, "Color Index", colorIndex,
			func(u values.Unifier) int { return u.Unified().(int) },
			func(int) string { return "%d" },
			-1, 255,
			func(newValue int) {
				view.requestPropertyChange(func(msg *text.ElectronicMessage) { msg.ColorIndex = newValue })
			})

		if view.model.currentKey.ID == ids.MailsStart {
			mailType := message.LeftDisplay - videoMailDisplayBase + 1
			selectedType := fmt.Sprintf("Unknown (0x%03X)", message.LeftDisplay)
			if mailType < 0 {
				mailType = 0
			}
			if mailType < len(knownMailTypes) {
				selectedType = knownMailTypes[mailType]
			}
			if readOnly {
				imgui.LabelText("Mail Type", selectedType)
			} else if imgui.BeginCombo("Mail Type", selectedType) {
				for typeIndex, typeString := range knownMailTypes {
					if imgui.SelectableV(typeString, mailType == typeIndex, 0, imgui.Vec2{}) {
						newIndex := -1
						if typeIndex > 0 {
							newIndex = videoMailDisplayBase + (typeIndex - 1)
						}
						view.requestPropertyChange(func(msg *text.ElectronicMessage) { msg.LeftDisplay = newIndex })
					}
				}
				imgui.EndCombo()
			}
		}
		if message.LeftDisplay < videoMailDisplayBase {
			view.renderDisplayGallery(readOnly, "Left Display", message.LeftDisplay, availableDisplays,
				func(newValue int) {
					view.requestPropertyChange(func(msg *text.ElectronicMessage) { msg.LeftDisplay = newValue })
				})
			view.renderDisplayGallery(readOnly, "Right Display", message.RightDisplay, availableDisplays,
				func(newValue int) {
					view.requestPropertyChange(func(msg *text.ElectronicMessage) { msg.RightDisplay = newValue })
				})
		}

		imgui.PopItemWidth()
	}
	imgui.EndChild()

	view.renderMFDs()
}

func (view *View) renderMFDs() {
	message, readOnly := view.currentMessage()

	view.renderSideImage("LeftMFD", message.LeftDisplay)
	imgui.SameLineV(0, 0)
	textToDisplay := message.VerboseText
	if imgui.BeginChildV("Text", imgui.Vec2{X: -250 * view.guiScale, Y: 190 * view.guiScale}, true, 0) {
		if !view.model.showVerboseText {
			textToDisplay = message.TerseText
		}

		imgui.PushTextWrapPos()
		if len(textToDisplay) == 0 {
			imgui.PushStyleColor(imgui.StyleColorText, imgui.Vec4{X: 1.0, Y: 1.0, Z: 1.0, W: 0.5})
			imgui.Text("(empty)")
			imgui.PopStyleColor()
		} else {
			imgui.Text(textToDisplay)
		}
		imgui.PopTextWrapPos()
	}
	imgui.EndChild()
	view.clipboardPopup(readOnly, "Text", textToDisplay, func(newValue string) {
		view.requestTextChange(func(msg *text.ElectronicMessage) {
			if view.model.showVerboseText {
				msg.VerboseText = newValue
			} else {
				msg.TerseText = newValue
			}
		})
	})
	imgui.SameLineV(0, 0)
	view.renderSideImage("RightMFD", message.RightDisplay)
}

func (view View) currentMessage() (msg text.ElectronicMessage, readOnly bool) {
	msg = view.messageOf(view.model.currentKey)
	readOnly = len(view.mod.ModifiedBlocks(view.model.currentKey.Lang, view.model.currentKey.ID.Plus(view.model.currentKey.Index))) == 0
	return
}

func (view View) hasAudio() bool {
	return (view.model.currentKey.ID == ids.MailsStart) || (view.model.currentKey.ID == ids.LogsStart)
}

func (view *View) currentSound() (sound audio.L8) {
	if view.hasAudio() {
		key := view.model.currentKey
		key.ID = key.ID.Plus(300)
		sound, _ = view.movieCache.Audio(key)
	}
	return
}

func (view View) messageOf(key resource.Key) text.ElectronicMessage {
	msg, cacheErr := view.messageCache.Message(key)
	if cacheErr != nil {
		msg = text.EmptyElectronicMessage()
	}
	return msg
}

func (view *View) availableDisplays() int {
	selector := view.mod.LocalizedResources(view.model.currentKey.Lang)
	res, err := selector.Select(ids.MfdDataBitmaps)
	if err != nil {
		return 0
	}
	return res.BlockCount()
}

func (view *View) renderText(readOnly bool, label string, value string, changeCallback func(string)) {
	imgui.LabelText(label, value)
	view.clipboardPopup(readOnly, label, value, changeCallback)
}

func (view *View) clipboardPopup(readOnly bool, label string, value string, changeCallback func(string)) {
	if imgui.BeginPopupContextItemV(label+"-Popup", 1) {
		if imgui.Selectable("Copy to Clipboard") {
			view.clipboard.SetString(value)
		}
		if !readOnly && imgui.Selectable("Copy from Clipboard") {
			newValue, err := view.clipboard.String()
			if err == nil {
				changeCallback(newValue)
			}
		}
		imgui.EndPopup()
	}
}

func (view *View) renderDisplayGallery(readOnly bool, label string, index, availableDisplays int,
	changeCallback func(int)) {
	unifier := values.NewUnifier()
	unifier.Add(index)
	values.RenderUnifiedSliderInt(readOnly, false, label, unifier,
		func(u values.Unifier) int { return u.Unified().(int) },
		func(int) string { return "%d" },
		-1, availableDisplays-1,
		changeCallback)
	render.TextureSelector("###"+label+" Texture", -1, view.guiScale, availableDisplays,
		index, view.imageCache,
		func(index int) resource.Key {
			return resource.KeyOf(ids.MfdDataBitmaps, view.model.currentKey.Lang, index)
		},
		func(index int) string { return fmt.Sprintf("%d", index) },
		changeCallback)
}

func (view *View) renderSideImage(label string, index int) {
	size := imgui.Vec2{X: 250 * view.guiScale, Y: 190 * view.guiScale}
	imgui.PushStyleVarVec2(imgui.StyleVarWindowPadding, imgui.Vec2{X: size.X, Y: -size.Y})
	if imgui.BeginChildV(label, size, false, 0) {
		if index >= 0 {
			key := resource.KeyOf(ids.MfdDataBitmaps, view.model.currentKey.Lang, index)
			render.TextureImage(label, view.imageCache, key, size)
		}
	}
	imgui.PopStyleVar()
	imgui.EndChild()
}

func (view *View) requestExportAudio(sound audio.L8) {
	filename := fmt.Sprintf("%05d_%s.wav",
		view.model.currentKey.ID.Plus(view.model.currentKey.Index).Plus(300).Value(),
		view.model.currentKey.Lang.String())

	external.ExportAudio(view.modalStateMachine, filename, sound)
}

func (view *View) requestImportAudio() {
	external.ImportAudio(view.modalStateMachine, func(sound audio.L8) {
		movieData := movie.ContainSoundData(sound)
		view.requestAudioChange(movieData)
	})
}

func (view *View) requestClear() {
	view.requestWipe(text.EmptyElectronicMessage().Encode(view.cp))
}

func (view *View) requestRemove() {
	view.requestWipe(nil)
}

func (view *View) requestWipe(newTextData [][]byte) {
	textEntries := make(map[resource.Language]messageDataEntry)
	audioEntries := make(map[resource.Language]messageDataEntry)
	textID := view.model.currentKey.ID.Plus(view.model.currentKey.Index)
	for lang := resource.Language(0); lang < resource.LanguageCount; lang++ {
		textEntries[lang] = messageDataEntry{
			oldData: view.mod.ModifiedBlocks(lang, textID),
			newData: newTextData,
		}
		if view.hasAudio() {
			audioEntries[lang] = messageDataEntry{
				oldData: view.mod.ModifiedBlocks(lang, textID.Plus(300)),
				newData: nil,
			}
		}
	}
	view.requestSetMessageData(textEntries, audioEntries)
}

func (view *View) requestTextChange(modifier func(*text.ElectronicMessage)) {
	msg := view.messageOf(view.model.currentKey)
	entries := make(map[resource.Language]messageDataEntry)
	modifier(&msg)

	entries[view.model.currentKey.Lang] = messageDataEntry{
		oldData: view.mod.ModifiedBlocks(view.model.currentKey.Lang, view.model.currentKey.ID.Plus(view.model.currentKey.Index)),
		newData: msg.Encode(view.cp),
	}
	view.requestSetMessageData(entries, nil)
}

func (view *View) requestPropertyChange(modifier func(*text.ElectronicMessage)) {
	entries := make(map[resource.Language]messageDataEntry)
	for lang := resource.Language(0); lang < resource.LanguageCount; lang++ {
		key := view.model.currentKey
		key.Lang = lang
		msg := view.messageOf(key)
		modifier(&msg)

		entries[lang] = messageDataEntry{
			oldData: view.mod.ModifiedBlocks(lang, key.ID.Plus(key.Index)),
			newData: msg.Encode(view.cp),
		}
	}
	view.requestSetMessageData(entries, nil)
}

func (view *View) requestAudioChange(movieData []byte) {
	entries := make(map[resource.Language]messageDataEntry)

	entries[view.model.currentKey.Lang] = messageDataEntry{
		oldData: view.mod.ModifiedBlocks(view.model.currentKey.Lang, view.model.currentKey.ID.Plus(view.model.currentKey.Index).Plus(300)),
		newData: [][]byte{movieData},
	}
	view.requestSetMessageData(nil, entries)
}

func (view *View) requestSetMessageData(textEntries map[resource.Language]messageDataEntry, audioEntries map[resource.Language]messageDataEntry) {
	command := setMessageDataCommand{
		key:             view.model.currentKey,
		showVerboseText: view.model.showVerboseText,
		model:           &view.model,
		textEntries:     textEntries,
		audioEntries:    audioEntries,
	}
	view.commander.Queue(command)
}
