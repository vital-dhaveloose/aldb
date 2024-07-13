In this text the basic ideas for ALDB are built from scratch.

- [ ] add visuals for each level, with labels to make them concrete examples
# Step 1: the statement
Information can exist in different forms. In this text we'll start with *symbolic* information, i.e. information represented using language. Symbolic information is the domain of letters, alphabets, words and sentences. Letters and words are mere building blocks; they don't convey information on their own. The minimal thing that does is a *statement* (also called a triplet): a combination of:
- a subject,
- a relation (also called property) and
- an object (sometimes called value or attribute).

For example: a document (subject) has a title (relation/property) "Progress report 2024Q3":

![[document-with-title.svg]]
# Step 2: grouping properties
Properties of a subject often come in groups. This can be because of different levels of abstraction of the subject each have their own set of properties, for example all images have a pixel size, but photo's have statements about the camera that was used while for digital drawings we might instead want to store some statements about which program was used to make them.
# Step 3: linking subjects
Objects of statements can be subjects themselves in other statements, or not. For example, the word count (relation/property) of a document (subject) can be 13 450 (object/value), a number that (for the intents and purposes of ALDB) doesn't have anything else going for it. But the document (object) can also be a description of (relation) a project (subject), that has lots of other relations.

When an object in one statement is also a subject on it's own, the statement can be seen as a link between these two subjects, with the relation labelling the link. When expanding this to multiple subjects and relations we see that symbolic information takes the form of a continuous, labelled graph, where each subject is related to each other subject via some route. Because of this, we will use the term "node" interchangeably with "subject" and "object". Statements thus become "links" in the graph. And so the "LD" in ALDB is explained: it stands for Linked Data and references the graph-like structure of symbolic information that is supported in the system.
	Linked Data is actually the name of a set of principles that ALDB is heavily inspired by, and means to implement.
# Step 4: add blobs
Not all information is symbolic though. We also use information systems to create, store and consume more direct input for the senses: images, sounds, video, etc. We call this *sensory* information.

> Blobs are pieces of data with unstructured contents or contents whose structure is not exposed for some reason (e.g. the structure being proprietary or the data being encrypted or compressed). As the goal of ALDB is to replace many different kinds of data storage, it must also be able to store, reference and manage blobs.

- [ ] not statements, but structure specific to the media type, can be abstracted as a blob -> ALD*B*
- [ ] we link this into our information graph by making it a type of object
# Step 5: activities for scoping
## Why we need scoping
The continuous, labelled graph we constructed in the previous steps is very useful, but some things appear difficult:
- Since there are many possible relations, and no way to distinguish important or usefulness between them, it's hard to achieve **overview** of a body of information. Think of a textbook with lots of cross references but without chapters. Overview help us construct mental models, take in new information, remember it and decrease the mental load while working with it, all by already providing the abstractions our minds need to function.
- **Navigating** to specific information that's not directly related to your starting point would not be trivial, again because the useful links would be hard to find between all other links. It would be like having to find a route to drive (without map or gps) when all street signs only point to the first upcoming intersection instead of cities etc. that are further away.
- Storing information requires **administering** it, for example setting up access control or synchronisation to devices for local access. Without a form of scoping this would be very tedious, as one would either have to do this on the level of individual items or by creating rules based on the many possible relations that link items together.

What's needed in these cases is a relation that provides scoping. One that expresses that one item is part of another. In ALDB it's called `is-part-of`.
## Structure
Typically, a strict hierarchy is chosen, i.e. a relation where each item has only one "parent" and there are no loops. While this constraint results in a number of useful features, it also has it's drawbacks, mainly that sometimes it's useful to put an item in multiple places in a structure. A typical use case is when multiple (groups of) people collaborate on a task or project. Each of them would have their own organisation structure in which it will have to be placed. Another use case is when different criteria are useful to organise, for example a chronological criterium like quarters; a departement/speciality criterium that distinguishes engineering, HR, finance, legal, ... Or sometimes an artifact like a table or a picture is used in multiple places (like a document and a slide deck).

That's why the `is-part-of` relation is not hierarchic, but forms a directed acyclic graph: nodes can have multiple "parents" but no cycles are allowed. Since cycles aren't allowed this implies that there will be roots, i.e. nodes that don't have a parent. Typically, these represent people's lives and organisations (see [[#Semantics]].

- [ ] if needed: sub relation is_mainly_part_of that does form a strict hierarchy
- [ ] "implying an "up" and a "down" direction."
- [ ] how is this useful for overview, navigation and administering?
## Semantics
And so we finally arrive at the first letter of ALDB: A stands for "activity structured". Activity is the name we give to all nodes that have "child" nodes in the `is-part-of` relation. With this relation they define a structure on all information in the system. In that way they resemble folders in a file system, but there are some important differences:
- The structure is not strictly hierarchical, as explained above.
- Activities are themselves nodes in the linked data graph, meaning they can have attributes and thus become richer descriptions of things than folders, that only have a name.

Another important difference is that folders can represent anything. They allow the user to organise using any criteria they want: by file type, chronologically, by importance (favourites), ... Activities on the other hand imply a certain meaning and organising criterium. They represent *things that a person or group spends time on*. More concrete examples are:
- projects, tasks, conversations, courses, conferences, social events, hobbies, vacations, etc.
- Groups of people are also often defined by an activity they cooperate on, like teams, departments or even whole organisations. In day-to-day language we wouldn't call these activities, but in ALDB they are represented as such.
- Another example of this is peoples lives: thes can also be seen as a grouping of activities. In fact, these are usually roots in the `is-part-of` relation.
- Artefacts such as documents, slide decks, datasets, … can become activities when they themselves contain tasks (like "write about ...") or conversations ("maybe we should change the order of the slides, because ..."). The artefact thus becomes synonymous with the work on the artefact.

Such semantics are more useful for the purposes listed above: overview, navigation and administering.
	A counter-example: on many operating systems each user has some folders that each represent a type of information: Documents, Music, Pictures, Video, etc. But this implies for example that the video's made during a vacation should be stored separately from the pictures that were taken, and that photo's that were made as part of some project should be stored separately from other project documentation like drawings and texts. Organising by information type is not useful for achieving overview (what information is available), navigation (e.g. back and forth between information types while working on a project) or administering information (e.g. sharing all project information with a new member).
- [ ] explain why
	- [ ] overview: 
	- [ ] navigation:
	- [ ] administering:
- [ ] Maybe also compare with organisation by place or time, by importance of preference, by aspect (technical, legal, economical, social, political, …), by step in a workflow (draft, active, archived, …), …
- [ ] reference [[2 Other data storage solutions]] for a more detailed comparison with hierarchic file systems
# Summary
ALDB supports:
- symbolic information, also called linked data (LD in ALDB) because of it's graph-like nature
- sensory information, implemented as blobs (B in ALDB)
- a useful scoping mechanism called activities (A in ALDB), with a flexible structure and guiding semantics
# Examples

- John Doe's life
	- Work
		- John Doe @ Green Corp
			- Odinson Offshore Wind Park
				- Business
					- Budget
					- Planning
					- Public tender
				- Engineering
					- Architecture
					- Mill selection
					- Placement optimisation
				- Environmental study
				- Legal
					- Permits
			- Hella Wind Park
				- Business
				- Engineering
					- Mill type selection
				- Legal
	- Doe family
		- Spain 2020
			- Itinerary
			- Playlist
			- Photos
		- Household
			- Groceries
		- Raspberry road 45
			- Technical
				- Electrical
				- Water
				- Information system
			- Business
				- Mortgage
				- Insurance
	- Business
		- Cash Flow
		- Insurance
		- Savings
	- Entertainment
		- Gaming
			- Halo
		- To watch/read