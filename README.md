# Elasticsearch CLI

A cli tool to interact with elasticsearch indices and perform elasticsearch operations. This project grew out of the necessities of my work to interact with elasticsearch. Mostly we try to follow the Elasticsearch API as documented in [Elasticsearch API](https://www.elastic.co/docs/api/doc/elasticsearch/) If you like to compile it, you would need the serde git repo in the same dir as the es-cli and preferably use a Go workspace. If you don't want to compile it, you will have to wait for a release of the binary.

Disclaimer: This project is WIP

## Working Assumptions

Currently the configuration assumes that the elasticsearch cluster is hosted on premise. Future versions could change that.

## Subcommands
  Currently the `es` command has a few subcommands. Some subcomannds have their own subcommand and most of subcommands has it's own flags.
  The current subcommands are:
  1. index. Has further subcommands 
    - create 
    - delete
    - list
    - clone

  2. count
  3. cluster
  4. search 

### Usage of Subcommands and Examples of flags
Below you can find some commands. For displaying the commands we will assume that the binary is in directory in the PATH of your user.
- Display helpful message: `es --help` or just `es`
- Retrieve an index's mappings:

  ```bash
  es index <index> --mappings
  ```

  or in shorthand notation

  ```bash
  es index <index> -m
  ```

- Create an index with provided settings and mappings. Usually I store the settings and mappings in a directory as json files, for example for the index <some-index> I have
  ```bash
  /elasticsearch/indices/some-index/mappings.json
  /elasticsearch/indices/some-index/settings.json
  ```
  To create an index we use the flag --directory (-d) with the path of the directory
  ```bash
  es index create --directory /elasticsearch/indices/some-index
  ```

- Delete an index

  ```bash
  es index delete <index>
  ```

- Doing simple searches on a single index, on multiple indices or using a wildcard

  ```bash
  es search <index>
  ```

- Doing search and retrieve only particular fields

  ```bash
  es search <index> --fields=field1,field2
  ```

To add more examples
