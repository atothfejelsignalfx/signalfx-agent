import re


def has_datapoint_with_metric_name(fake_services, metric_name):
    for dp in fake_services.datapoints:
        if dp.metric == metric_name:
            return True
    return False


# Tests if any datapoint received has the given dim key/value on it.
def has_datapoint_with_dim(fake_services, key, value):
    for dp in fake_services.datapoints:
        for dim in dp.dimensions:
            if dim.key == key and dim.value == value:
                return True
    return False

# Tests if a command run against a container returns with an exit code of 0
def container_cmd_exit_0(container, command):
    code, _ = container.exec_run(command)
    return code == 0


# This won't work very robustly if the text spans multiple lines.
def text_is_in_stream(stream, text):
    return text.encode("utf-8") in b"".join(stream.readlines())


def has_log_message(output, level="info", message=""):
    for l in output.splitlines():
        m = re.search(r'(?<=level=)\w+', l)
        if m is None:
            continue
        if level == m.group(0) and message in l:
            return True
    return False
