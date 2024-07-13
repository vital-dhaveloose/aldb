import { Participation } from './participation'

export interface ActivityVersion {
    /**
     * Schema meta field to reference this schema with.
     */
    $schema?: string
    /**
     * The identifier of the Activity, formatted as a URI.
     */
    id?: string
    /**
     * A version representation of the ActivityVersion
     */
    version?: string
    /**
     * A map with short strings for representing the Activity in a UI. The keys of the map are locales.
     */
    label?: {
      [k: string]: string
    }
    /**
     * The period in which the Activity is considered 'current'.
     */
    period?: {
      startTime?: string
      endTime?: string
    }
    participations?: Participation[]
    /**
     * Activities that are part of this Activity.
     */
    subs?: ActivityVersion[]
    /**
     * Activities that this Activity is part of.
     */
    supers?: ActivityVersion[]
    attributeSets?: {
      [k: string]: {
        manifest?: AttributeSetManifest
        attributes?: {
          [k: string]: unknown
        }
      }
    }
    blob?: {
      manifest?: BlobManifest
      bytesBase64?: string
      blobRef?: string
    }
  }

  export interface AttributeSetManifest {
    id?: string
  }

  export interface BlobManifest {
    mediaType?: string
  }
  