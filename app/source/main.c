#include <stdint.h>
#include <stdio.h>

#define SDL_MAIN_USE_CALLBACKS 1
#include <SDL3/SDL.h>
#include <SDL3/SDL_main.h>

#include "hacked/nuklear/NuklearSdlBridge.h"

struct HackEdApp
{
   SDL_Window *window;
   SDL_Renderer *renderer;
   struct nk_context *ctx;
};

static float appGetBaseUIScale(SDL_Window *const window)
{
   (void)window;
   return 1.0f; // SDL_GetWindowDisplayScale(window) * 1.0f;
}

static SDL_AppResult appFailSDL(char const *const message)
{
   SDL_LogError(SDL_LOG_CATEGORY_CUSTOM, "Error; %s: %s", message, SDL_GetError());
   return SDL_APP_FAILURE;
}

static struct nk_color color(uint8_t const r, uint8_t const g, uint8_t const b, float const a)
{
   struct nk_color const value = {.r = r, .g = g, .b = b, .a = (uint8_t)SDL_roundf(255.0f * a)};
   return value;
}

static struct nk_color colorDoubleFull(float const a)
{
   return color(0xC4, 0x38, 0x9F, a);
}

static struct nk_color colorDoubleDark(float const a)
{
   return color(0x31, 0x01, 0x38, a);
}

static struct nk_color colorTripleFull(float const a)
{
   return color(0x21, 0xFF, 0x43, a);
}

static struct nk_color colorTripleDark(float const a)
{
   return color(0x06, 0xCC, 0x94, a);
}

static struct nk_color colorTripleLight(float const a)
{
   return color(0x51, 0x99, 0x58, a);
}

static void appSetUIStyle(struct nk_context *const ctx)
{
   struct nk_vec2 const zero = nk_vec2(0.0f, 0.0f);

   struct nk_color hackedColorStyle[NK_COLOR_COUNT] = {0};
   for (size_t i = 0; i < NK_COLOR_COUNT; i++)
   {
      hackedColorStyle[i] = color(0xFF, 0x00, 0x00, 0x75f);
   }
   hackedColorStyle[NK_COLOR_TEXT] = color(0x5B, 0xAC, 0x1E, 1.0f);
   hackedColorStyle[NK_COLOR_WINDOW] = colorDoubleDark(1.0f);
   hackedColorStyle[NK_COLOR_HEADER] = colorTripleLight(0.70f);
   hackedColorStyle[NK_COLOR_BORDER] = colorDoubleFull(1.0f);

   hackedColorStyle[NK_COLOR_BUTTON] = colorTripleDark(0.4f);
   hackedColorStyle[NK_COLOR_BUTTON_HOVER] = colorTripleFull(1.0f);
   hackedColorStyle[NK_COLOR_BUTTON_ACTIVE] = colorDoubleFull(1.0f);

   hackedColorStyle[NK_COLOR_TOGGLE] = colorTripleDark(0.4f);
   hackedColorStyle[NK_COLOR_TOGGLE_HOVER] = colorTripleFull(1.0f);
   hackedColorStyle[NK_COLOR_TOGGLE_CURSOR] = colorDoubleFull(1.0f);

   hackedColorStyle[NK_COLOR_SELECT] = colorTripleDark(0.4f);
   hackedColorStyle[NK_COLOR_SELECT_ACTIVE] = colorDoubleFull(1.0f);

   hackedColorStyle[NK_COLOR_SLIDER] = colorTripleDark(0.4f);
   hackedColorStyle[NK_COLOR_SLIDER_CURSOR] = colorTripleLight(1.0f);
   hackedColorStyle[NK_COLOR_SLIDER_CURSOR_HOVER] = colorTripleFull(1.0f);
   hackedColorStyle[NK_COLOR_SLIDER_CURSOR_ACTIVE] = colorDoubleFull(1.0f);

   hackedColorStyle[NK_COLOR_PROPERTY] = colorDoubleDark(1.0f);

   hackedColorStyle[NK_COLOR_EDIT] = colorDoubleDark(1.0f);
   hackedColorStyle[NK_COLOR_EDIT_CURSOR] = colorDoubleFull(1.0f);

   hackedColorStyle[NK_COLOR_COMBO] = colorDoubleDark(1.0f);

   hackedColorStyle[NK_COLOR_CHART] = colorDoubleDark(1.0f);
   hackedColorStyle[NK_COLOR_CHART_COLOR] = colorTripleFull(1.0f);
   hackedColorStyle[NK_COLOR_CHART_COLOR_HIGHLIGHT] = colorTripleLight(1.0f);

   hackedColorStyle[NK_COLOR_SCROLLBAR] = colorTripleDark(1.0f);
   hackedColorStyle[NK_COLOR_SCROLLBAR_CURSOR] = colorTripleLight(1.0f);
   hackedColorStyle[NK_COLOR_SCROLLBAR_CURSOR_HOVER] = colorTripleFull(1.0f);
   hackedColorStyle[NK_COLOR_SCROLLBAR_CURSOR_ACTIVE] = colorDoubleFull(1.0f);

   hackedColorStyle[NK_COLOR_TAB_HEADER] = colorTripleLight(0.54f);

   hackedColorStyle[NK_COLOR_KNOB] = colorTripleDark(1.0f);
   hackedColorStyle[NK_COLOR_KNOB_CURSOR] = colorTripleLight(1.0f);
   hackedColorStyle[NK_COLOR_KNOB_CURSOR_HOVER] = colorTripleFull(1.0f);
   hackedColorStyle[NK_COLOR_KNOB_CURSOR_ACTIVE] = colorDoubleFull(1.0f);

   nk_style_from_table(ctx, hackedColorStyle);

   struct nk_style *style = &ctx->style;
   /* default text */
   struct nk_style_text *text = &style->text;
   // text->color = table[NK_COLOR_TEXT];
   text->padding = nk_vec2(0, 0);
   text->color_factor = 1.0f;
   text->disabled_factor = NK_WIDGET_DISABLED_FACTOR;

   /* default button */
   struct nk_style_button *button = &style->button;
   /*
    button->normal                     = nk_style_item_color(table[NK_COLOR_BUTTON]);
    button->hover                      = nk_style_item_color(table[NK_COLOR_BUTTON_HOVER]);
    button->active                     = nk_style_item_color(table[NK_COLOR_BUTTON_ACTIVE]);
    button->border_color               = table[NK_COLOR_BORDER];
    button->text_background            = table[NK_COLOR_BUTTON];
    button->text_normal                = table[NK_COLOR_TEXT];
    button->text_hover                 = table[NK_COLOR_TEXT];
    button->text_active                = table[NK_COLOR_TEXT];
   */
   button->padding = zero; // nk_vec2(2.0f, 2.0f);
   button->image_padding = zero; // nk_vec2(0.0f, 0.0f);
   button->touch_padding = zero; // nk_vec2(0.0f, 0.0f);
   button->userdata = nk_handle_ptr(NULL);
   button->text_alignment = NK_TEXT_CENTERED;
   button->border = 1.0f;
   button->rounding = 0.0f;
   button->color_factor_text = 1.0f;
   button->color_factor_background = 1.0f;
   button->disabled_factor = NK_WIDGET_DISABLED_FACTOR;

   /* contextual button */
   button = &style->contextual_button;
   /*
    button->normal          = nk_style_item_color(table[NK_COLOR_WINDOW]);
    button->hover           = nk_style_item_color(table[NK_COLOR_BUTTON_HOVER]);
    button->active          = nk_style_item_color(table[NK_COLOR_BUTTON_ACTIVE]);
    button->border_color    = table[NK_COLOR_WINDOW];
    button->text_background = table[NK_COLOR_WINDOW];
    button->text_normal     = table[NK_COLOR_TEXT];
    button->text_hover      = table[NK_COLOR_TEXT];
    button->text_active     = table[NK_COLOR_TEXT];
   */
   button->padding = zero; //(2.0f, 2.0f);
   button->touch_padding = zero; //(0.0f, 0.0f);
   button->userdata = nk_handle_ptr(NULL);
   button->text_alignment = NK_TEXT_CENTERED;
   button->border = 0.0f;
   button->rounding = 0.0f;
   button->color_factor_text = 1.0f;
   button->color_factor_background = 1.0f;
   button->disabled_factor = NK_WIDGET_DISABLED_FACTOR;

   /* menu button */
   button = &style->menu_button;
   /*
   button->normal = nk_style_item_color(table[NK_COLOR_WINDOW]);
   button->hover = nk_style_item_color(table[NK_COLOR_WINDOW]);
   button->active = nk_style_item_color(table[NK_COLOR_WINDOW]);
   button->border_color = table[NK_COLOR_WINDOW];
   button->text_background = table[NK_COLOR_WINDOW];
   button->text_normal = table[NK_COLOR_TEXT];
   button->text_hover = table[NK_COLOR_TEXT];
   button->text_active = table[NK_COLOR_TEXT];
   */
   button->padding = zero; //(2.0f, 2.0f);
   button->touch_padding = zero; //(0.0f, 0.0f);
   button->text_alignment = NK_TEXT_CENTERED;
   button->border = 0.0f;
   button->rounding = 0.0f;
   button->color_factor_text = 1.0f;
   button->color_factor_background = 1.0f;
   button->disabled_factor = NK_WIDGET_DISABLED_FACTOR;
#if 0
   /* checkbox toggle */
   struct nk_style_toggle *toggle = &style->checkbox;
   /*
   toggle->normal = nk_style_item_color(table[NK_COLOR_TOGGLE]);
   toggle->hover = nk_style_item_color(table[NK_COLOR_TOGGLE_HOVER]);
   toggle->active = nk_style_item_color(table[NK_COLOR_TOGGLE_HOVER]);
   toggle->cursor_normal = nk_style_item_color(table[NK_COLOR_TOGGLE_CURSOR]);
   toggle->cursor_hover = nk_style_item_color(table[NK_COLOR_TOGGLE_CURSOR]);
   toggle->userdata = nk_handle_ptr(0);
   toggle->text_background = table[NK_COLOR_WINDOW];
   toggle->text_normal = table[NK_COLOR_TEXT];
   toggle->text_hover = table[NK_COLOR_TEXT];
   toggle->text_active = table[NK_COLOR_TEXT];
   */
   toggle->padding = zero; //(2.0f, 2.0f);
   toggle->touch_padding = zero; //(0, 0);
   toggle->border_color = nk_rgba(0, 0, 0, 0);
   toggle->border = 0.0f;
   toggle->spacing = 4;
   toggle->color_factor = 1.0f;
   toggle->disabled_factor = NK_WIDGET_DISABLED_FACTOR;

   /* option toggle */
   toggle = &style->option;
   /*
   toggle->normal = nk_style_item_color(table[NK_COLOR_TOGGLE]);
   toggle->hover = nk_style_item_color(table[NK_COLOR_TOGGLE_HOVER]);
   toggle->active = nk_style_item_color(table[NK_COLOR_TOGGLE_HOVER]);
   toggle->cursor_normal = nk_style_item_color(table[NK_COLOR_TOGGLE_CURSOR]);
   toggle->cursor_hover = nk_style_item_color(table[NK_COLOR_TOGGLE_CURSOR]);
   toggle->userdata = nk_handle_ptr(0);
   toggle->text_background = table[NK_COLOR_WINDOW];
   toggle->text_normal = table[NK_COLOR_TEXT];
   toggle->text_hover = table[NK_COLOR_TEXT];
   toggle->text_active = table[NK_COLOR_TEXT];
   */
   toggle->padding = zero; //(3.0f, 3.0f);
   toggle->touch_padding = zero; //(0, 0);
   toggle->border_color = nk_rgba(0, 0, 0, 0);
   toggle->border = 0.0f;
   toggle->spacing = 4;
   toggle->color_factor = 1.0f;
   toggle->disabled_factor = NK_WIDGET_DISABLED_FACTOR;

   /* selectable */
   struct nk_style_selectable *select = &style->selectable;
   /*
   select->normal = nk_style_item_color(table[NK_COLOR_SELECT]);
   select->hover = nk_style_item_color(table[NK_COLOR_SELECT]);
   select->pressed = nk_style_item_color(table[NK_COLOR_SELECT]);
   select->normal_active = nk_style_item_color(table[NK_COLOR_SELECT_ACTIVE]);
   select->hover_active = nk_style_item_color(table[NK_COLOR_SELECT_ACTIVE]);
   select->pressed_active = nk_style_item_color(table[NK_COLOR_SELECT_ACTIVE]);
   select->text_normal = table[NK_COLOR_TEXT];
   select->text_hover = table[NK_COLOR_TEXT];
   select->text_pressed = table[NK_COLOR_TEXT];
   select->text_normal_active = table[NK_COLOR_TEXT];
   select->text_hover_active = table[NK_COLOR_TEXT];
   select->text_pressed_active = table[NK_COLOR_TEXT];
   */
   select->padding = zero; //(2.0f, 2.0f);
   select->image_padding = zero; //(2.0f, 2.0f);
   select->touch_padding = zero; //(0, 0);
   select->rounding = 0.0f;
   select->color_factor = 1.0f;
   select->disabled_factor = NK_WIDGET_DISABLED_FACTOR;

   /* slider */
   struct nk_style_slider *slider = &style->slider;
   slider->normal = nk_style_item_hide();
   slider->hover = nk_style_item_hide();
   slider->active = nk_style_item_hide();
   /*
   slider->bar_normal = table[NK_COLOR_SLIDER];
   slider->bar_hover = table[NK_COLOR_SLIDER];
   slider->bar_active = table[NK_COLOR_SLIDER];
   slider->bar_filled = table[NK_COLOR_SLIDER_CURSOR];
   slider->cursor_normal = nk_style_item_color(table[NK_COLOR_SLIDER_CURSOR]);
   slider->cursor_hover = nk_style_item_color(table[NK_COLOR_SLIDER_CURSOR_HOVER]);
   slider->cursor_active = nk_style_item_color(table[NK_COLOR_SLIDER_CURSOR_ACTIVE]);
   */
   slider->inc_symbol = NK_SYMBOL_TRIANGLE_RIGHT;
   slider->dec_symbol = NK_SYMBOL_TRIANGLE_LEFT;
   slider->cursor_size = zero; //(16, 16);
   slider->padding = zero; //(2, 2);
   slider->spacing = zero; //(2, 2);
   slider->show_buttons = nk_false;
   slider->bar_height = 4;
   slider->rounding = 0;
   slider->color_factor = 1.0f;
   slider->disabled_factor = NK_WIDGET_DISABLED_FACTOR;

   /* slider buttons */
   button = &style->slider.inc_button;
   /*
   button->normal = nk_style_item_color(nk_rgb(40, 40, 40));
   button->hover = nk_style_item_color(nk_rgb(42, 42, 42));
   button->active = nk_style_item_color(nk_rgb(44, 44, 44));
   button->border_color = nk_rgb(65, 65, 65);
   button->text_background = nk_rgb(40, 40, 40);
   button->text_normal = nk_rgb(175, 175, 175);
   button->text_hover = nk_rgb(175, 175, 175);
   button->text_active = nk_rgb(175, 175, 175);
   */
   button->padding = zero; //(8.0f, 8.0f);
   button->touch_padding = zero; //(0.0f, 0.0f);
   button->text_alignment = NK_TEXT_CENTERED;
   button->border = 1.0f;
   button->rounding = 0.0f;
   button->color_factor_text = 1.0f;
   button->color_factor_background = 1.0f;
   button->disabled_factor = NK_WIDGET_DISABLED_FACTOR;
   style->slider.dec_button = style->slider.inc_button;

   /* knob */
   struct nk_style_knob *knob = &style->knob;
   knob->normal = nk_style_item_hide();
   knob->hover = nk_style_item_hide();
   knob->active = nk_style_item_hide();
   /*
   knob->knob_normal = table[NK_COLOR_KNOB];
   knob->knob_hover = table[NK_COLOR_KNOB];
   knob->knob_active = table[NK_COLOR_KNOB];
   knob->cursor_normal = table[NK_COLOR_KNOB_CURSOR];
   knob->cursor_hover = table[NK_COLOR_KNOB_CURSOR_HOVER];
   knob->cursor_active = table[NK_COLOR_KNOB_CURSOR_ACTIVE];

   knob->knob_border_color = table[NK_COLOR_BORDER];
   */
   knob->knob_border = 1.0f;

   knob->padding = zero; //(2, 2);
   knob->spacing = zero; //(2, 2);
   knob->cursor_width = 2;
   knob->color_factor = 1.0f;
   knob->disabled_factor = NK_WIDGET_DISABLED_FACTOR;

   /* progressbar */
   struct nk_style_progress *prog = &style->progress;
   /*
   prog->normal = nk_style_item_color(table[NK_COLOR_SLIDER]);
   prog->hover = nk_style_item_color(table[NK_COLOR_SLIDER]);
   prog->active = nk_style_item_color(table[NK_COLOR_SLIDER]);
   prog->cursor_normal = nk_style_item_color(table[NK_COLOR_SLIDER_CURSOR]);
   prog->cursor_hover = nk_style_item_color(table[NK_COLOR_SLIDER_CURSOR_HOVER]);
   prog->cursor_active = nk_style_item_color(table[NK_COLOR_SLIDER_CURSOR_ACTIVE]);
   prog->border_color = nk_rgba(0, 0, 0, 0);
   prog->cursor_border_color = nk_rgba(0, 0, 0, 0);
   */
   prog->padding = zero; //(4, 4);
   prog->rounding = 0;
   prog->border = 0;
   prog->cursor_rounding = 0;
   prog->cursor_border = 0;
   prog->color_factor = 1.0f;
   prog->disabled_factor = NK_WIDGET_DISABLED_FACTOR;

   /* scrollbars */
   struct nk_style_scrollbar *scroll = &style->scrollh;
   /*
   scroll->normal = nk_style_item_color(table[NK_COLOR_SCROLLBAR]);
   scroll->hover = nk_style_item_color(table[NK_COLOR_SCROLLBAR]);
   scroll->active = nk_style_item_color(table[NK_COLOR_SCROLLBAR]);
   scroll->cursor_normal = nk_style_item_color(table[NK_COLOR_SCROLLBAR_CURSOR]);
   scroll->cursor_hover = nk_style_item_color(table[NK_COLOR_SCROLLBAR_CURSOR_HOVER]);
   scroll->cursor_active = nk_style_item_color(table[NK_COLOR_SCROLLBAR_CURSOR_ACTIVE]);
   */
   scroll->dec_symbol = NK_SYMBOL_CIRCLE_SOLID;
   scroll->inc_symbol = NK_SYMBOL_CIRCLE_SOLID;
   /*
   scroll->border_color = table[NK_COLOR_SCROLLBAR];
   scroll->cursor_border_color = table[NK_COLOR_SCROLLBAR];
   */
   scroll->padding = zero; //(0, 0);
   scroll->show_buttons = nk_false;
   scroll->border = 0;
   scroll->rounding = 0;
   scroll->border_cursor = 0;
   scroll->rounding_cursor = 0;
   scroll->color_factor = 1.0f;
   scroll->disabled_factor = NK_WIDGET_DISABLED_FACTOR;
   style->scrollv = style->scrollh;

   /* scrollbars buttons */
   button = &style->scrollh.inc_button;
   /*
   button->normal = nk_style_item_color(nk_rgb(40, 40, 40));
   button->hover = nk_style_item_color(nk_rgb(42, 42, 42));
   button->active = nk_style_item_color(nk_rgb(44, 44, 44));
   button->border_color = nk_rgb(65, 65, 65);
   button->text_background = nk_rgb(40, 40, 40);
   button->text_normal = nk_rgb(175, 175, 175);
   button->text_hover = nk_rgb(175, 175, 175);
   button->text_active = nk_rgb(175, 175, 175);
   */
   button->padding = zero; //(4.0f, 4.0f);
   button->touch_padding = zero; //(0.0f, 0.0f);
   button->text_alignment = NK_TEXT_CENTERED;
   button->border = 1.0f;
   button->rounding = 0.0f;
   button->color_factor_text = 1.0f;
   button->color_factor_background = 1.0f;
   button->disabled_factor = NK_WIDGET_DISABLED_FACTOR;
   style->scrollh.dec_button = style->scrollh.inc_button;
   style->scrollv.inc_button = style->scrollh.inc_button;
   style->scrollv.dec_button = style->scrollh.inc_button;

   /* edit */
   struct nk_style_edit *edit = &style->edit;
   /*
   edit->normal = nk_style_item_color(table[NK_COLOR_EDIT]);
   edit->hover = nk_style_item_color(table[NK_COLOR_EDIT]);
   edit->active = nk_style_item_color(table[NK_COLOR_EDIT]);
   edit->cursor_normal = table[NK_COLOR_TEXT];
   edit->cursor_hover = table[NK_COLOR_TEXT];
   edit->cursor_text_normal = table[NK_COLOR_EDIT];
   edit->cursor_text_hover = table[NK_COLOR_EDIT];
   edit->border_color = table[NK_COLOR_BORDER];
   edit->text_normal = table[NK_COLOR_TEXT];
   edit->text_hover = table[NK_COLOR_TEXT];
   edit->text_active = table[NK_COLOR_TEXT];
   edit->selected_normal = table[NK_COLOR_TEXT];
   edit->selected_hover = table[NK_COLOR_TEXT];
   edit->selected_text_normal = table[NK_COLOR_EDIT];
   edit->selected_text_hover = table[NK_COLOR_EDIT];
   */
   edit->scrollbar_size = zero; //(10, 10);
   edit->scrollbar = style->scrollv;
   edit->padding = zero; //(4, 4);
   edit->row_padding = 2;
   edit->cursor_size = 4;
   edit->border = 1;
   edit->rounding = 0;
   edit->color_factor = 1.0f;
   edit->disabled_factor = NK_WIDGET_DISABLED_FACTOR;

   /* property */
   struct nk_style_property *property = &style->property;
   /*
   property->normal = nk_style_item_color(table[NK_COLOR_PROPERTY]);
   property->hover = nk_style_item_color(table[NK_COLOR_PROPERTY]);
   property->active = nk_style_item_color(table[NK_COLOR_PROPERTY]);
   property->border_color = table[NK_COLOR_BORDER];
   property->label_normal = table[NK_COLOR_TEXT];
   property->label_hover = table[NK_COLOR_TEXT];
   property->label_active = table[NK_COLOR_TEXT];
   */
   property->sym_left = NK_SYMBOL_TRIANGLE_LEFT;
   property->sym_right = NK_SYMBOL_TRIANGLE_RIGHT;
   property->padding = zero; //(4, 4);
   property->border = 1;
   property->rounding = 0;
   property->color_factor = 1.0f;
   property->disabled_factor = NK_WIDGET_DISABLED_FACTOR;

   /* property buttons */
   button = &style->property.dec_button;
   /*
   button->normal = nk_style_item_color(table[NK_COLOR_PROPERTY]);
   button->hover = nk_style_item_color(table[NK_COLOR_PROPERTY]);
   button->active = nk_style_item_color(table[NK_COLOR_PROPERTY]);
   button->border_color = nk_rgba(0, 0, 0, 0);
   button->text_background = table[NK_COLOR_PROPERTY];
   button->text_normal = table[NK_COLOR_TEXT];
   button->text_hover = table[NK_COLOR_TEXT];
   button->text_active = table[NK_COLOR_TEXT];
   */
   button->padding = zero; //(0.0f, 0.0f);
   button->touch_padding = zero; //(0.0f, 0.0f);
   button->text_alignment = NK_TEXT_CENTERED;
   button->border = 0.0f;
   button->rounding = 0.0f;
   button->color_factor_text = 1.0f;
   button->color_factor_background = 1.0f;
   button->disabled_factor = NK_WIDGET_DISABLED_FACTOR;
   style->property.inc_button = style->property.dec_button;

   /* property edit */
   edit = &style->property.edit;
   /*
   edit->normal = nk_style_item_color(table[NK_COLOR_PROPERTY]);
   edit->hover = nk_style_item_color(table[NK_COLOR_PROPERTY]);
   edit->active = nk_style_item_color(table[NK_COLOR_PROPERTY]);
   edit->border_color = nk_rgba(0, 0, 0, 0);
   edit->cursor_normal = table[NK_COLOR_TEXT];
   edit->cursor_hover = table[NK_COLOR_TEXT];
   edit->cursor_text_normal = table[NK_COLOR_EDIT];
   edit->cursor_text_hover = table[NK_COLOR_EDIT];
   edit->text_normal = table[NK_COLOR_TEXT];
   edit->text_hover = table[NK_COLOR_TEXT];
   edit->text_active = table[NK_COLOR_TEXT];
   edit->selected_normal = table[NK_COLOR_TEXT];
   edit->selected_hover = table[NK_COLOR_TEXT];
   edit->selected_text_normal = table[NK_COLOR_EDIT];
   edit->selected_text_hover = table[NK_COLOR_EDIT];
   */
   edit->padding = zero; //(0, 0);
   edit->cursor_size = 8;
   edit->border = 0;
   edit->rounding = 0;
   edit->color_factor = 1.0f;
   edit->disabled_factor = NK_WIDGET_DISABLED_FACTOR;

   /* chart */
   struct nk_style_chart *chart = &style->chart;
   /*
   chart->background = nk_style_item_color(table[NK_COLOR_CHART]);
   chart->border_color = table[NK_COLOR_BORDER];
   chart->selected_color = table[NK_COLOR_CHART_COLOR_HIGHLIGHT];
   chart->color = table[NK_COLOR_CHART_COLOR];
   */
   chart->padding = zero; //(4, 4);
   chart->border = 0;
   chart->rounding = 0;
   chart->color_factor = 1.0f;
   chart->disabled_factor = NK_WIDGET_DISABLED_FACTOR;
   chart->show_markers = nk_true;

   /* combo */
   struct nk_style_combo *combo = &style->combo;
   /*
   combo->normal = nk_style_item_color(table[NK_COLOR_COMBO]);
   combo->hover = nk_style_item_color(table[NK_COLOR_COMBO]);
   combo->active = nk_style_item_color(table[NK_COLOR_COMBO]);
   combo->border_color = table[NK_COLOR_BORDER];
   combo->label_normal = table[NK_COLOR_TEXT];
   combo->label_hover = table[NK_COLOR_TEXT];
   combo->label_active = table[NK_COLOR_TEXT];
   */
   combo->sym_normal = NK_SYMBOL_TRIANGLE_DOWN;
   combo->sym_hover = NK_SYMBOL_TRIANGLE_DOWN;
   combo->sym_active = NK_SYMBOL_TRIANGLE_DOWN;
   combo->content_padding = zero; //(4, 4);
   combo->button_padding = zero; //(0, 4);
   combo->spacing = zero; //(4, 0);
   combo->border = 1;
   combo->rounding = 0;
   combo->color_factor = 1.0f;
   combo->disabled_factor = NK_WIDGET_DISABLED_FACTOR;

   /* combo button */
   button = &style->combo.button;
   /*
   button->normal = nk_style_item_color(table[NK_COLOR_COMBO]);
   button->hover = nk_style_item_color(table[NK_COLOR_COMBO]);
   button->active = nk_style_item_color(table[NK_COLOR_COMBO]);
   button->border_color = nk_rgba(0, 0, 0, 0);
   button->text_background = table[NK_COLOR_COMBO];
   button->text_normal = table[NK_COLOR_TEXT];
   button->text_hover = table[NK_COLOR_TEXT];
   button->text_active = table[NK_COLOR_TEXT];
   */
   button->padding = zero; //(2.0f, 2.0f);
   button->touch_padding = zero; //(0.0f, 0.0f);
   button->text_alignment = NK_TEXT_CENTERED;
   button->border = 0.0f;
   button->rounding = 0.0f;
   button->color_factor_text = 1.0f;
   button->color_factor_background = 1.0f;
   button->disabled_factor = NK_WIDGET_DISABLED_FACTOR;

   /* tab */
   struct nk_style_tab *tab = &style->tab;
   /*
   tab->background = nk_style_item_color(table[NK_COLOR_TAB_HEADER]);
   tab->border_color = table[NK_COLOR_BORDER];
   tab->text = table[NK_COLOR_TEXT];
   */
   tab->sym_minimize = NK_SYMBOL_TRIANGLE_RIGHT;
   tab->sym_maximize = NK_SYMBOL_TRIANGLE_DOWN;
   tab->padding = zero; //(4, 4);
   tab->spacing = zero; //(4, 4);
   tab->indent = 10.0f;
   tab->border = 1;
   tab->rounding = 0;
   tab->color_factor = 1.0f;
   tab->disabled_factor = NK_WIDGET_DISABLED_FACTOR;

   /* tab button */
   button = &style->tab.tab_minimize_button;
   /*
   button->normal = nk_style_item_color(table[NK_COLOR_TAB_HEADER]);
   button->hover = nk_style_item_color(table[NK_COLOR_TAB_HEADER]);
   button->active = nk_style_item_color(table[NK_COLOR_TAB_HEADER]);
   button->border_color = nk_rgba(0, 0, 0, 0);
   button->text_background = table[NK_COLOR_TAB_HEADER];
   button->text_normal = table[NK_COLOR_TEXT];
   button->text_hover = table[NK_COLOR_TEXT];
   button->text_active = table[NK_COLOR_TEXT];
   */
   button->padding = zero; //(2.0f, 2.0f);
   button->touch_padding = zero; //(0.0f, 0.0f);
   button->text_alignment = NK_TEXT_CENTERED;
   button->border = 0.0f;
   button->rounding = 0.0f;
   button->color_factor_text = 1.0f;
   button->color_factor_background = 1.0f;
   button->disabled_factor = NK_WIDGET_DISABLED_FACTOR;
   style->tab.tab_maximize_button = *button;

   /* node button */
   button = &style->tab.node_minimize_button;
   /*
   button->normal = nk_style_item_color(table[NK_COLOR_WINDOW]);
   button->hover = nk_style_item_color(table[NK_COLOR_WINDOW]);
   button->active = nk_style_item_color(table[NK_COLOR_WINDOW]);
   button->border_color = nk_rgba(0, 0, 0, 0);
   button->text_background = table[NK_COLOR_TAB_HEADER];
   button->text_normal = table[NK_COLOR_TEXT];
   button->text_hover = table[NK_COLOR_TEXT];
   button->text_active = table[NK_COLOR_TEXT];
   */
   button->padding = zero; //(2.0f, 2.0f);
   button->touch_padding = zero; //(0.0f, 0.0f);
   button->text_alignment = NK_TEXT_CENTERED;
   button->border = 0.0f;
   button->rounding = 0.0f;
   button->color_factor_text = 1.0f;
   button->color_factor_background = 1.0f;
   button->disabled_factor = NK_WIDGET_DISABLED_FACTOR;
   style->tab.node_maximize_button = *button;
#endif
   /* window header */
   struct nk_style_window *win = &style->window;
   win->header.align = NK_HEADER_RIGHT;
   win->header.close_symbol = NK_SYMBOL_X;
   win->header.minimize_symbol = NK_SYMBOL_MINUS;
   win->header.maximize_symbol = NK_SYMBOL_PLUS;
   /*
   win->header.normal = nk_style_item_color(table[NK_COLOR_HEADER]);
   win->header.hover = nk_style_item_color(table[NK_COLOR_HEADER]);
   win->header.active = nk_style_item_color(table[NK_COLOR_HEADER]);
   win->header.label_normal = table[NK_COLOR_TEXT];
   win->header.label_hover = table[NK_COLOR_TEXT];
   win->header.label_active = table[NK_COLOR_TEXT];
   */
   win->header.label_padding = zero; //(4, 4);
   win->header.padding = zero; //(4, 4);
   win->header.spacing = zero; //(0, 0);

   /* window header close button */
   button = &style->window.header.close_button;
   /*
   button->normal = nk_style_item_color(table[NK_COLOR_HEADER]);
   button->hover = nk_style_item_color(table[NK_COLOR_HEADER]);
   button->active = nk_style_item_color(table[NK_COLOR_HEADER]);
   button->border_color = nk_rgba(0, 0, 0, 0);
   button->text_background = table[NK_COLOR_HEADER];
   button->text_normal = table[NK_COLOR_TEXT];
   button->text_hover = table[NK_COLOR_TEXT];
   button->text_active = table[NK_COLOR_TEXT];
   */
   button->padding = zero; //(0.0f, 0.0f);
   button->touch_padding = zero; //(0.0f, 0.0f);
   button->text_alignment = NK_TEXT_CENTERED;
   button->border = 0.0f;
   button->rounding = 0.0f;
   button->color_factor_text = 1.0f;
   button->color_factor_background = 1.0f;
   button->disabled_factor = NK_WIDGET_DISABLED_FACTOR;

   /* window header minimize button */
   button = &style->window.header.minimize_button;
   /*
   button->normal = nk_style_item_color(table[NK_COLOR_HEADER]);
   button->hover = nk_style_item_color(table[NK_COLOR_HEADER]);
   button->active = nk_style_item_color(table[NK_COLOR_HEADER]);
   button->border_color = nk_rgba(0, 0, 0, 0);
   button->text_background = table[NK_COLOR_HEADER];
   button->text_normal = table[NK_COLOR_TEXT];
   button->text_hover = table[NK_COLOR_TEXT];
   button->text_active = table[NK_COLOR_TEXT];
   */
   button->padding = zero; //(0.0f, 0.0f);
   button->touch_padding = zero; //(0.0f, 0.0f);
   button->text_alignment = NK_TEXT_CENTERED;
   button->border = 0.0f;
   button->rounding = 0.0f;
   button->color_factor_text = 1.0f;
   button->color_factor_background = 1.0f;
   button->disabled_factor = NK_WIDGET_DISABLED_FACTOR;

   /* window */
   /*
   win->background = table[NK_COLOR_WINDOW];
   win->fixed_background = nk_style_item_color(table[NK_COLOR_WINDOW]);
   win->border_color = table[NK_COLOR_BORDER];
   win->popup_border_color = table[NK_COLOR_BORDER];
   win->combo_border_color = table[NK_COLOR_BORDER];
   win->contextual_border_color = table[NK_COLOR_BORDER];
   win->menu_border_color = table[NK_COLOR_BORDER];
   win->group_border_color = table[NK_COLOR_BORDER];
   win->tooltip_border_color = table[NK_COLOR_BORDER];
   win->scaler = nk_style_item_color(table[NK_COLOR_TEXT]);
   */

   win->rounding = 0.0f;
   win->spacing = zero; //(4, 4);
   win->scrollbar_size = nk_vec2(5, 5);
   win->min_size = nk_vec2(6, 6); //(64, 64);

   win->combo_border = 1.0f;
   win->contextual_border = 0.0f;
   win->menu_border = 0.0f;
   win->group_border = 0.0f;
   win->tooltip_border = 1.0f;
   win->popup_border = 0.0f;
   win->border = 1.0f;
   win->min_row_height_padding = 0;

   win->padding = zero; // set to zero because it affects menu bars
   win->group_padding = nk_vec2(1, 1); //(4, 4);
   win->popup_padding = nk_vec2(1, 1); //(4, 4);
   win->combo_padding = nk_vec2(1, 1); //(4, 4);
   win->contextual_padding = nk_vec2(1, 1); //(4, 4);
   win->menu_padding = nk_vec2(1, 1); //(4, 4);
   win->tooltip_padding = nk_vec2(1, 1); //(4, 4);

   win->tooltip_origin = NK_TOP_LEFT;
   win->tooltip_offset = zero; //(12, 12);
   win->tooltip_delay = 0.5f;
}

SDL_AppResult SDL_AppInit(void **const appstate, int const argc, char *argv[])
{
   (void)argc;
   (void)argv;

   if (!SDL_Init(SDL_INIT_VIDEO | SDL_INIT_EVENTS))
   {
      return appFailSDL("failed to initialize SDL");
   }

   SDL_SetAppMetadata("InkyBlackness - HackEd", REPO_SHORT_VERSION, "io.github.inkyblackness.hacked");
   if (!SDL_Init(SDL_INIT_VIDEO | SDL_INIT_EVENTS))
   {
      return appFailSDL("failed to initialize SDL");
   }
   struct HackEdApp *app = SDL_malloc(sizeof(*app));
   if (app == NULL)
   {
      return appFailSDL("failed to allocate application memory");
   }
   if (!SDL_CreateWindowAndRenderer(
          "InkyBlackness - HackEd - " REPO_LONG_VERSION, 320, 200, SDL_WINDOW_RESIZABLE | SDL_WINDOW_HIGH_PIXEL_DENSITY, &app->window, &app->renderer))
   {
      SDL_free(app);
      return appFailSDL("failed to create window/renderer");
   }
   *appstate = app;

   if (!SDL_SetRenderVSync(app->renderer, 1))
   {
      SDL_LogError(SDL_LOG_CATEGORY_CUSTOM, "SDL_SetRenderVSync failed: %s", SDL_GetError());
   }
   SDL_SetRenderLogicalPresentation(app->renderer, 320, 200, SDL_LOGICAL_PRESENTATION_LETTERBOX);
   SDL_SetRenderDrawBlendMode(app->renderer, SDL_BLENDMODE_BLEND); // Ensure blend mode is set on all platforms

   struct nk_context *ctx = uiBridgeInit(app->window, app->renderer);
   app->ctx = ctx;

   {
      float const scale = appGetBaseUIScale(app->window);
      uiBridgeSetFont(ctx, scale);

      appSetUIStyle(ctx);
   }
   nk_input_begin(ctx);

   return SDL_APP_CONTINUE;
}

SDL_AppResult SDL_AppEvent(void *const appstate, SDL_Event *const event)
{
   struct HackEdApp const *const app = appstate;

   switch (event->type)
   {
   case SDL_EVENT_QUIT:
      return SDL_APP_SUCCESS;
   case SDL_EVENT_KEY_DOWN:
      if (event->key.key == SDLK_Q && event->key.mod & SDL_KMOD_CTRL)
      {
         return SDL_APP_SUCCESS;
      }
      break;
   case SDL_EVENT_WINDOW_DISPLAY_SCALE_CHANGED:
      /* You may wish to rescale the renderer and Nuklear during this event.
       * Without this the UI and Font could appear too small or too big.
       * This is not handled by the demo in order to keep it simple,
       * but you may wish to re-bake the Font whenever this happens. */
      SDL_Log("Unhandled scale event! Nuklear may appear blurry");
      return SDL_APP_CONTINUE;
   default:
      break;
   }

   SDL_ConvertEventToRenderCoordinates(app->renderer, event);
   uiBridgeHandleEvent(app->ctx, event);

   return SDL_APP_CONTINUE;
}

static void appClearBackground(SDL_Renderer *const renderer)
{
   SDL_SetRenderDrawColor(renderer, 0x00, 0x00, 0x00, 0xFF);
   SDL_RenderClear(renderer);
}

SDL_AppResult SDL_AppIterate(void *const appstate)
{
   struct HackEdApp const *const app = appstate;
   struct nk_context *const ctx = app->ctx;
   SDL_AppResult appResult = SDL_APP_CONTINUE;
   float const scale = appGetBaseUIScale(app->window);

   static uint64_t previous = 0;

   uint64_t now = SDL_GetTicksNS();
   if (previous == 0)
   {
      previous = now;
   }
   uint64_t elapsed = now - previous;
   if (elapsed == 0)
   {
      elapsed = 1;
   }
   double const fps = 1000000000.0 / (double)elapsed;
   char fpsLine[30];
   sprintf(fpsLine, "%3lumsec - %.1f", (unsigned long)(elapsed / 1000000ULL), fps);
   previous = now;
   nk_input_end(ctx);

   if (nk_begin(ctx, "main-menu", nk_rect(0, 0, 10000, 8 * scale), NK_WINDOW_NO_SCROLLBAR))
   {
      nk_menubar_begin(ctx);
      {
         nk_layout_row_static(ctx, 8 * scale, 200, 2);
         if (nk_menu_begin_label(ctx, "File", 0, nk_vec2(100, 100)))
         {
            nk_layout_row_dynamic(ctx, 10 * scale, 1);
            if (nk_menu_item_label(ctx, "New...", 0))
            {
               SDL_Log("new...");
            }
            if (nk_menu_item_label(ctx, "Quit", 0))
            {
               appResult = SDL_APP_SUCCESS;
            }
            nk_menu_end(ctx);
         }
         if (nk_menu_begin_label(ctx, fpsLine, 0, nk_vec2(100, 100)))
         {
            nk_layout_row_dynamic(ctx, 10 * scale, 1);
            nk_menu_end(ctx);
         }
      }
      nk_menubar_end(ctx);
   }
   nk_end(ctx);
   /* GUI */
   if (nk_begin(ctx, "Demo", nk_rect(50, 50, 230, 190), NK_WINDOW_BORDER | NK_WINDOW_MOVABLE | NK_WINDOW_SCALABLE | NK_WINDOW_MINIMIZABLE | NK_WINDOW_TITLE))
   {
      enum
      {
         EASY,
         HARD
      };
      static int op = EASY;
      static int property = 20;
      enum
      {
         textBufferSize = 1000
      };
      static char textBuffer[textBufferSize] = "";
      static int textBufferUsedLen = 0;

      nk_layout_row_static(ctx, 9, 60, 1);
      if (nk_button_label(ctx, "button"))
      {
         SDL_Log("button pressed");
      }
      nk_layout_row_dynamic(ctx, 30, 2);
      if (nk_option_label(ctx, "easy", op == EASY))
         op = EASY;
      if (nk_option_label(ctx, "hard", op == HARD))
         op = HARD;
      nk_layout_row_dynamic(ctx, 25, 1);
      nk_property_int(ctx, "Compression:", 0, &property, 1000, 1, 1);

      nk_layout_row_dynamic(ctx, 20, 1);
      nk_label(ctx, "background:", NK_TEXT_LEFT);
      nk_layout_row_dynamic(ctx, 25, 1);

      nk_edit_string(ctx, NK_EDIT_ALWAYS_INSERT_MODE | NK_EDIT_MULTILINE, textBuffer, &textBufferUsedLen, textBufferSize, NULL);
   }
   nk_end(ctx);

   appClearBackground(app->renderer);

   uiBridgeRender(ctx);
   uiBridgeUpdateTextInput(ctx);

   SDL_RenderPresent(app->renderer);

   nk_input_begin(ctx);
   return appResult;
}

void SDL_AppQuit(void *const appstate, SDL_AppResult const result)
{
   (void)result;

   struct HackEdApp *const app = appstate;
   if (app != NULL)
   {
      SDL_Log("Quitting");
      nk_input_end(app->ctx);
      uiBridgeShutdown(app->ctx);
      SDL_DestroyRenderer(app->renderer);
      SDL_DestroyWindow(app->window);
      SDL_free(app);
   }
}

char *nk_sdl_dtoa(char *const str, double const d)
{
   NK_ASSERT(str);
   if (!str)
      return NULL;
   (void)SDL_snprintf(str, 99999, "%.17g", d);
   return str;
}
