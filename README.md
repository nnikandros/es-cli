# Elasticsearch CLI

A cli tool to interact with elasticsearch indices and perform elasticsearch operations. This project grew out of the necessities of my work to interact with elasticsearch.

## Working Assumptions

Currently the configuration assumes that the elasticsearch cluster is hosted on premise. Future versions could change that.

## Example of Commands

- Display help `es --help`
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

  ```bash
  es index create --directory /elasticsearch/indices/new-index
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
