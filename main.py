import yaml
import os
import subprocess
from flask import Flask, request
from jinja2 import Template

base_dir = os.path.dirname(__file__)
app = Flask(__name__)

def get_content(fname):
    with open(fname) as f:
        return f.read()

def get_yaml(fname):
    return yaml.safe_load(get_content(fname))

def get_pipeline(req):
    """
    redirect to pipeline file name
    base on definition in https://docs.drone.io/extensions/configuration/
    then load pipeline and return
    """
    index = get_yaml("index.yml")
    slug = req["repo"]["slug"]
    if slug in index["pipeline"]:
        return get_yaml(f"pipelines/{index['pipeline'][slug]}.yml")

def get_template(pipeline):
    """
    get template base on pipeline config
    """
    return get_yaml(f"templates/{pipeline['template']}.yml")

def render_pipeline(pipeline, template):
    """
    load module in template
    render module by pipeline values
    return render results
    """
    result = []
    for module in template:
        result.append(yaml.safe_load(Template(get_content(f"modules/{module}.j2")).render(pipeline[module] if module in pipeline else {})))
    return result

@app.route("/", methods=["POST"])
def run():
    req = request.json
    pipeline = get_pipeline(req)
    template = get_template(pipeline)
    result = render_pipeline(pipeline, template)
    if len(result) > 0:
        return {"data": yaml.safe_dump_all(result)}
    else:
        return "", 204

if __name__ == "__main__":
    host = os.environ.get("HOST", "0.0.0.0")
    port = os.environ.get("PORT", "3000")
    app.run(host=host, port=port, debug=True)