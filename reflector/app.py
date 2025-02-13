import logging
import db
from flask import Flask, jsonify, render_template, request, redirect, url_for
import subprocess
import os

logging.basicConfig(
    level=logging.INFO,
)
app = Flask("reflector")


@app.route("/", methods=["GET"])
def first_page():
    return render_template("login.html", content="Reflector UI")


@app.route("/landing_page", methods=["POST"])
def landing_page():
    email = request.form["email"]
    return render_template("landing_page.html", content=email)


@app.route("/echo", methods=["GET"])
def echo():
    key = request.args.get("key", "")
    return f"<html><body>echo: {key}</body></html>"


@app.route("/sqlinjection", methods=["GET", "POST"])
def index():
    search_query = request.form["search"]
    results = db.get_comments(search_query)
    return render_template(
        "sql_injection.html", results=results, search_query=search_query
    )


# New route using Bootstrap's JavaScript from a CDN
@app.route("/external-script", methods=["GET"])
def external_script():
    message = request.args.get("message", "")

    # Include an external script from Bootstrap's CDN
    script = '<script src="https://cdn.jsdelivr.net/npm/bootstrap@5.1.3/dist/js/bootstrap.bundle.min.js"></script>'

    # Your HTML structure with Bootstrap JS integration
    html = f"""
        <html>
        <head>{script}</head>
        <body>
            <div class="container">
                <h1>Bootstrap Example</h1>
                <p>{message}</p>
                <button id="alertButton" class="btn btn-primary">Click Me!</button>
            </div>
            <script>
                // Using Bootstrap's JavaScript functionality
                document.getElementById("alertButton").addEventListener("click", function() {{
                    alert("Button clicked! Message: {message}");
                }});
            </script>
        </body>
        </html>
    """

    return html


@app.route("/shellinjection", methods=["GET", "POST"])
def shellinjection():
    cwd = os.getcwd()
    search_query = request.form["data"]
    result = subprocess.run(
        "sh " + cwd + "/shell.sh " + search_query,
        capture_output=True,
        shell=True,
    )
    return render_template(
        "shell_injection.html", results=[result.stdout], search_query=search_query
    )
# @app.route('/')
# def upload_form():
#     return render_template('upload.html')

@app.route('/upload_file',methods=["GET", "POST"])
def upload_file():
    # Check if the post request has the file part
    if 'file' not in request.files:
        return 'No file part'
    UPLOAD_FOLDER = os.getcwd()
    file = request.files['file']
    
    # If the user does not select a file
    if file.filename == '':
        return 'No selected file'
    
    if file :
        filename = file.filename
        filepath = os.path.join(UPLOAD_FOLDER, filename)
        file.save(filepath)
        return f'File successfully uploaded to {filepath}'
    
    return 'Invalid file type'

# Route for the 'redeem.html' page
@app.route('/redeem',methods=["POST"])
def redeem():
    return render_template('redeem.html')

if __name__ == "__main__":
    app.run(host="0.0.0.0", port=8080, debug=True)
