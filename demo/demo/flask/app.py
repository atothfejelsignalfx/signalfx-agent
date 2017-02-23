from flask import Flask

import signalfx
import os

endpoint = 'http://{0}:{1}'.format(
    os.environ['SIGNALFX_METRIC_PROXY_SERVICE_HOST'], os.environ['SIGNALFX_METRIC_PROXY_SERVICE_PORT'])

app = Flask(__name__)

counter = 0

@app.route('/')
def hello_world():
    global counter

    with signalfx.SignalFx().ingest(token='', endpoint=endpoint) as sfx:
        sfx.send(
            cumulative_counters=[dict(
                metric='users.visit',
                value=counter,
            )])
        counter += 1
    return 'Hello, World!'

app.run(host='0.0.0.0')
