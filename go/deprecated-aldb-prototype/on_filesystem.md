# API implementation working on top of a file system

ActivityId's <-> file and folder names
VersionId's <-> ?
Structured data <-> file contents (yaml? json? rdf? others? support multiple? default?)
Blob data <-> file contents
Primary relation <-> containers in this relation get a folder


e.g.

green-corp.yml
green-corp/
    odinson-offshore-wind-park.yml
    odinson-offshore-wind-park/
        engineering.yml
        engineering/
            architecture/
            mill-selection/
            placement-optimisation/
        environmental-study/
        legal/

Folders <-> primary link/relation (is_primary_part_of/is_primary_container_of)
Files   <-> blobs and structured data  

Bypassing API (= direct changes in de file system)?
    - can break links