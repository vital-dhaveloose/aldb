
  export interface Participation {
    participator?: Person | Group
    /**
     * An array of URIs referencing the roles of the participators in this participation.
     */
    roles?: string[]
    [k: string]: unknown
  }

  export interface Person {
    givenName?: string
    familyName?: string
    email?: string
    [k: string]: unknown
  }
  
  /**
   * A group of people.
   */
  export interface Group {
    /**
     * A string describing the Group, for example: "Team Blue"
     */
    display?: string
    email?: string
    [k: string]: unknown
  }