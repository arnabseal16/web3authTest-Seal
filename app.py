from flask import Flask, jsonify, request

app = Flask(__name__)

@app.route('/')
def hello_word():
    return jsonify(
        {
            'greetings_from': request.remote_addr
        }
    )

if __name__ == '__main__':
    app.run(host='0.0.0.0', port=8000)