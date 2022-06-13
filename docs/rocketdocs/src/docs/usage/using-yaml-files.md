---
title: Using YAML files
description: Learn how to YAML files in Rocket Docs.
---

To use YAML files, create a new YAML file in a `yamlFiles` folder and it will be exposed as `all{FileName}Yaml` GraphQL node.

To store YAML in another folder, please change `yamlFilesPath` option of `@rocketseat/gatsby-theme-docs` plugin which is configured in `gatsby.config.js`.

## Example file

If our project has `letters.yml` or `letters.yaml` under a `yamlFiles` folder which looks like:

```yml
- character: a
- character: b
- character: c
```

Then the following three nodes would be created:

```json
[
  {
    "character": "a"
  },
  {
    "character": "b"
  },
  {
    "character": "c"
  }
]
```

## Query GraphQL node

We can query the GraphQL nodes which is created from the YAML file as the following code:

```graphql
{
  allLettersYaml {
    edges {
      node {
        character
      }
    }
  }
}
```

Which would return:

```json
{
  "allLettersYaml": {
    "edges": [
      {
        "node": {
          "character": "a"
        }
      },
      {
        "node": {
          "character": "b"
        }
      },
      {
        "node": {
          "character": "c"
        }
      }
    ]
  }
}
```

More information, please refer to [gatsby-transformer-yaml plugin page](https://www.gatsbyjs.com/plugins/gatsby-transformer-yaml/).
