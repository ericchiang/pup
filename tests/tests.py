
from subprocess import Popen, PIPE, STDOUT

example_data = """
<html>
  <head>
  </head>
  <body>
    <div>
      <div class="nav clearfix">
        My data
      </div>
      <p>Some other data</p>
    </div>
  </body>
</html>
"""

# run a pup command as a subprocess
def run_pup(args, input_html):
    cmd = ["pup"]
    cmd.extend(args)
    p = Popen(cmd, stdout=PIPE, stdin=PIPE, stderr=PIPE)
    stdout_data = p.communicate(input=input_html)[0]
    p.wait()
    return stdout_data

# simply count the number of lines returned by this pup command
def run_pup_count(args, input_html):
    pup_output = run_pup(args, input_html)
    lines = [l for l in pup_output.split("\n") if l]
    return len(lines)

def test_class_selector():
    assert run_pup_count([".nav"], example_data) == 3

def test_attr_eq():
    assert run_pup_count(["[class=nav]"], example_data) == 0

def test_attr_pre():
    assert run_pup_count(["[class^=nav]"], example_data) == 3
    assert run_pup_count(["[class^=clearfix]"], example_data) == 0

def test_attr_post():
    assert run_pup_count(["[class$=nav]"], example_data) == 0
    assert run_pup_count(["[class$=clearfix]"], example_data) == 3

def test_attr_func():
    result = run_pup(["div", "attr{class}"], example_data).strip()
    assert result == ""
    result = run_pup(["div", "div", "attr{class}"], example_data).strip()
    assert result == "nav clearfix"

def test_text_func():
    result = run_pup(["p", "text{}"], example_data).strip()
    assert result == "Some other data"
