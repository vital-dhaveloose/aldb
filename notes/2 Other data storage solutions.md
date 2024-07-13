In this text some popular data storage solutions are compared to ALDB.

# Aspects being compared
- the reasons a special relation "is part of" exists in ALDB: the need to get overview, to navigate (starting from scratch) and administering data
- support for different kinds of information in the structure
	- symbolic/sensory
		- on the leaves and/or on the containers in the structure
		- universal addressing
	- formal/informal: notes as flat text or markup, support for tables, etc.
		- [ ] **Why are notes special amongst other kinds of information?**
- constraints on the structure:
	- number of levels
	- number of items in a group
	- separation between containers and contents
	- multiple relations (= multiple ways to group)
	- multiple locations within one relation's structure?
		- different people might want different locations/structure (easier to remember structure if chosen yourself --> overview and navigation)
	- semantics of the relation: expectations/guidance?
- [FAIR](https://en.wikipedia.org/wiki/FAIR_data) aspects
	- Findable, Accessible, Interoperable, Reusable
- shared structure across apps (define and manage once, refine and use everywhere)
- global solution or multiple roots per person
- connectedness (related to value) vs. fragmentation
- ownership and privacy
- freedom of choosing applications and other products (<> vendor lock in)

# Hierarchic file systems
Hierarchic file systems are characterised by:
- their hierarchic layout and the natural consequence that each item has one location
- the separation between containers (directories/folders) and contents (files)

- [ ] Local vs. cloud

- [ ] overview
- [ ] navigation
	- [ ] from scratch: from some root
	- [ ] direct: path/url for direct navigation
		- [ ] no difference between location in context and identity
- [ ] data administration
	- [ ] built on hierarchic scoping, generally works well
- [ ] kinds of information
	- [ ] no good support for structured/symbolic data (limited schema available, schema almost never part of the file system)
	- [ ] no built-in standardisation for informal information (notes)
- [ ] structure constraints
	- [ ] no practical limits on number of levels
	- [ ] everyone must share one single structure
	- [ ] no guiding principle for semantics
- [ ] FAIR
	- [ ] Findable
	- [ ] Accessible: API?
	- [ ] Interoperable: open schema?
	- [ ] Reusable
# Specialist apps
Data stored only (or at least mainly) to be accessed through a single app. Not only contents, but structure is managed by the app. For example Outlook: folders in the app but only one file in the file system. Or allmost all mobile apps and web apps (Cloud based apps): data storage is completely hidden from the user, structure and data management is custom made (e.g. Facebook, Microsoft To Do, WhatsApp, ...)

Specialized for a specific type of information, for example tasks, finances, 

- [ ] overview
- [ ] navigation
	- [ ] from/to other apps dependent on the app ("deep links" or "intents")
- [ ] data management
	- [ ] custom (for access control: re-model users and groups for each app, authentication can also become fragmented)
- [ ] kinds of information
	- [ ] specific to app functionality
	- [ ] sometimes little extensions, such as a "notes" field or support for attachment or comments/conversations
- [ ] structure constraints
	- [ ] generally quite limited structure (e.g. limited number of levels, each with their own semantics)
	- [ ] generally no shared structure, so repeat structures in every app
		- [ ] different structural constraints in each app makes this extra difficult (e.g. map a multi-level hierarchy to a flat list or a hierarchy with limited levels)
# Generalist apps
Exactly like specialist apps, but support information in many forms, such as (richt) text, user defined tables, images and other multimedia, included files (attachments), canvasses, etc.

- [ ] OneNote
- [ ] Notion
- [ ] Jira
- [ ] Spreadsheets (and other tabular data with user defined schema, like Access)
- [ ] Nano
- [ ] Canvas apps like Miro, Mural
- [ ] Slack
	- [ ] not only communication, but also files, notes in canvasses, etc.

# Obsidian
Like a generalist app, but based on the file system and very open.


# Relational databases
Not as a back-end storage technology that apps are built on, but more as a system that users directly interact with, for example by writing queries or constructing views (like pivot tables in Excel).

# Linked Data (W3)
- [ ] URI's for naming and referencing things, HTTP for CRUD operations, RDF for structured data
	- See also [Linked Data Platform (LDP)](https://www.w3.org/TR/ldp/), [Introduction to: Linked Data Platform](https://www.dataversity.net/introduction-linked-data-platform/) and [Learning W3C Linked Data Platform with examples](https://www.slideshare.net/nandana/learning-w3c-linked-data-platform-with-examples)

- [ ] Solid

- [ ] overview
- [ ] navigation
- [ ] data management
- [ ] kinds of information
- [ ] structure constraints
# (Social) Semantic Desktops
- [ ] Why they "failed" and why ALDB doesn't make the same mistakes
- [ ] See also OneNote page "Studie (Social) Semantic Desktop"
# ALDB
- [ ] discuss all aspects for ALDB
# Overview and conclusion: why ALDB?