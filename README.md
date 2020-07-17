# Drone Pipelines

Drone multiple pipeline management base on Flask and Jinja.

## Features

- Simple http server as configuration extension
- Management: modules, templates, and actual pipelines
- Template syntax base on Jinja

## Run

```bash
python main.py
```

and pass address to drone's `DRONE_YAML_ENDPOINT`, see [Configuration Extension](https://docs.drone.io/extensions/configuration/)

## Management

1. Define your own modules with Jinja
2. Compose modules to template, see `templates` directory
3. Write pipeline: choose one template, pass variable to modules