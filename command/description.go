package command

const (
	GIT_PROGRAM_SHORT_DESC = "The stupid content tracker"
	GIT_PROGRAM_LONG_DESC  = `Git is a fast, scalable, distributed revision control system with an unusually rich command set that provides both high-level operations and full access to internals.`

	INIT_COMMIT_SHORT_DESC = "Create an empty Git repository or reinitialize an existing one"
	INIT_COMMIT_LONG_DESC  = `This command creates an empty Git repository - basically a .git directory with subdirectories for objects, refs/heads, refs/tags, and template files. An initial branch without any commits will be created.`

	ADD_COMMIT_SHORT_DESC = "Add file contents to the index"
	ADD_COMMIT_LONG_DESC  = `This command updates the index using the current content found in the working tree, to prepare the content staged for the next commit. It typically adds the current content of existing paths as a whole, but with some options it can also be used to add content with only part of the changes made to the working tree files applied, or remove paths that do not exist in the working tree anymore.`

	COMMIT_COMMIT_SHORT_DESC = "Record changes to the repository"
	COMMIT_COMMIT_LONG_DESC  = `Stores the current contents of the index in a new commit along with a log message from the user describing the changes.`
)
