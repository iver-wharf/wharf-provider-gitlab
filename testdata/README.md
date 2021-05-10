# Git Lab importer test data

Files list:

- `projects_all.json` - all projects. This file is splitted into three next
  files:

  - `groups/default_9/projects.json` - projects from group with ID 9 (Default)

  - `groups/default_9/super-project_84/projects.json.json` - projects from
    group with ID 84 (Default / Super-project)

  - `groups/basket_25/mushroom_87/projects.json` - projects from group with
    ID 87 (Basket / Mushroom)

## System hierarchy

Below how the data looks like. In parenthesis 'P' when project and 'G' when
group as prefix before the ID.

```text
-> Default (G9) 
  -> Super-project (G84)
    -> web (P267)
    -> builder (P252)
    -> docs (P225)
  -> main_test-proj (P277)
  -> super-project-messages (P261)

-> Basket (G25)
  -> Mushroom (G87)
  -> Boletus(P264)
```

Total projects count: 6.
