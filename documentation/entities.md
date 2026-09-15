# Data layers and storage paths

Data layers is the general term for any resource-contributing sources. Typically, the base layer would be the game's
original resources, set as a standard entry in the editor's settings.

The final layer on the other side of the list is the storage for the loose resources of the project.

## Entities and components

```mermaid
erDiagram
    EditorSettings {
    }
    EditorSettingsTag {
    }
    EditorSettings |o--o| EditorSettingsTag: has
    EditorSettings |o--o{ DataLayer: LoadsResourcesFrom

    LoadsResourcesFrom {
        uint8_t order
    }
    class LoadsResourcesFrom component

    Project {
    }
    ProjectTag {
    }
    Project |o--o| ProjectTag: has
    Project |o--o| StoragePath: has
    Project |o--o{ DataLayer: LoadsResourcesFrom

    DataLayer {
    }
    DataLayer |o--o| DataLayerPath: has

    StoragePath {
        String path
    }
    class StoragePath component

    DataLayerPath {
        String path
    }
    class DataLayerPath component

    classDef component fill:#0FC
```

# Resources

Resources are handled by the `ResourceSystem`, based on the `DataLayerPath`s, as well as what the editor overrides.

* Resource entities have a `ResourceId` component.
  They are created by either things that require them, or the `ResourceSystem` to indicate their availability.
* If a resource is edited as part of the current project, they will point to the corresponding resource data via
  'EditedAs'.
* If something requires a resource, the `ResourceSystem` will perform a load and then associate resource data
  via 'LoadedAs'.
* `ResourceData` entities then have components according to their types. Media data will have their dynamically
  allocated data as part of the components. Detailed data structures will have other components.

## Entities and components

```mermaid
erDiagram
    Resource {
    }
    Resource |o--o| ResourceId: has
    Resource |o--o| ResourceData: LoadedAs
    Resource |o--o| ResourceData: EditedAs
    Something }o--o{ Resource: Requires

    ResourceId {
        uint16_t id
        uint16_t index
    }
    class ResourceId component

    LoadedAs {
    }
    class LoadedAs component
    EditedAs {
    }
    class EditedAs component


    ResourceData {
    }
    ResourceData |o--o| FontData: has
    ResourceData |o--o| BitmapData: has

    FontData {
        Font data
    }
    class FontData component

    BitmapData {
        Bitmap data
    }
    class BitmapData component

    Texture {
    }
    Texture |o--o| TextureId: has
    Texture |o--o| TextureProperties: LoadedAs
    Texture |o--o| TextureProperties: EditedAs
    Something }o--o{ Texture: Requires
    TextureId {
        uint16_t id
    }
    class TextureId component

    Object {
    }
    Object |o--o| ObjectTriple: has
    Object |o--o| ObjectProperties: LoadedAs
    Object |o--o| ObjectProperties: EditedAs
    Something }o--o{ Object: Requires
    ObjectTriple {
        uint8_t class
        uint8_t subclass
        uint8_t type
    }
    class ObjectTriple component

    classDef component fill: #0FC
```

