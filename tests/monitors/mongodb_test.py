from functools import partial as p
import os
import string

from tests.helpers import fake_backend
from tests.helpers.util import wait_for, run_agent, run_container
from tests.helpers.assertions import *

mongo_config = string.Template("""
monitors:
  - type: collectd/mongodb
    host: $host
    port: 27017
    databases: [admin]
""")

def test_mongo():
    with run_container("mongo:3.6") as mongo_cont:
        config = mongo_config.substitute(host=mongo_cont.attrs["NetworkSettings"]["IPAddress"])

        with run_agent(config) as [backend, _, _]:
            assert wait_for(p(has_datapoint_with_dim, backend, "plugin", "mongo")), "Didn't get mongo datapoints"

