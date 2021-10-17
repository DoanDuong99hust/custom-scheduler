import flask
import random
from flask import request, jsonify

app = flask.Flask(__name__)
app.config["DEBUG"] = True

random.seed(1)
rand1 = random.randint(1,100)
random.seed(2)
rand2 = random.randint(1,100)
random.seed(3)
rand3 = random.randint(1,100)
randomNum =random.randint(1,100)
  
@app.route('/api/v1/test', methods=['GET'])
def cpu_data():
    data = {
    "switch": [
        {
            "port": [
                {
                    "Tx": randomNum,
                    "Rx": randomNum,
                },
                {
                    "Tx": rand1,
                    "Rx": rand1
                },
                {
                    "Tx": rand2,
                    "Rx": rand2
                }
            ]
        },
        {
            "port": [
                {
                    "Tx": rand1,
                    "Rx": rand1
                },
                {
                    "Tx": randomNum,
                    "Rx": randomNum
                },
                {
                    "Tx": rand3,
                    "Rx": rand3
                }
            ]
        },
        {
            "port": [
                {
                    "Tx": rand2,
                    "Rx": rand2
                },
                {
                    "Tx": rand3,
                    "Rx": rand3
                },
                {
                    "Tx": randomNum,
                    "Rx": randomNum
                }
            ]
        }
    ]
}

    return jsonify(data)

app.run()