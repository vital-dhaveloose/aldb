Identificatie en selectie
- alle, geen, specifieke (één, meerdere, alle behalve), van-tot
  ordebepaling voor van-tot?
  (lineairiseerbaarheid van versies impliceert niet dat elke versie "afstamt" van de vorige versie, branching is nog altijd mogelijk)

- obv index
    first <-> #0
    latest <-> #-1
    
    alternatief: intuitief voor wie van 1 nummert (ipv vanaf 0) + goeie default
        kan genoemd worden: versionNumber (<> versionIndex)
        previous <-> #-1
        latest <-> #0 (= latest = current = default)
        first <-> #1
        second <-> #2

- mix obv index en id, bv. {foo, #-1} <-> foo en de recentste

- kruisproduct activities / activity versions met '@'
    - bv.: `(id-a, id-b)@(#-2, #-1)` <-> laatste twee versies van a en van b

- meerdere selecties in één: unie, doorsnee

- obv relaties (soort van GraphQL)

- omzetting van/naar string (of json of yaml)
     *          <--> All
     {}         <--> None
     foo        <--> Only "foo"
     {foo,bar}  <--> "foo" and "bar"
     !{foo,bar} <--> Not "foo" or "bar"
     [foo,bar[  <--> between "foo" and "bar"
     [foo,[     <--> greater than "foo"
     #0         <--> first
     #-1        <--> latest
     ^regex$    <--> accepted by the regex

compatibiliteit met:
- URI/URL (https://en.wikipedia.org/wiki/Uniform_Resource_Identifier)
    mss iets anders dan '@' als separator tussen activityId en versie? bv. '|', '/'

    URI = scheme ":" ["//" [userinfo "@"] host [":" port]] path ["?" query] ["#" fragment]
    URN = "urn:" nis ":" nss [["?+" r-component] ["?=" q-component]] ["#" fragment]

- UUIDv4/Base64
- RFC3996 (datetime) bv. 2022-10-30T19:46:00.000Z
- regex
- Base64
- hex notation