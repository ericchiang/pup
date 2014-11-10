# Tests

A simple set of tests to help maintain sanity.

These tests don't actually test functionality they only make sure pup behaves
the same after code changes.

`cmds.txt` holds a list of commands to perform on `index.html`.

The output of each of these commands produces a specific sha1sum. The expected
sha1sum of each command is in `expected_output.txt`.

Running the `test` file (just a bash script) will run the tests and diff the
output. If pup has changed at all since the last version, you'll see the sha1sums
that changed and the commands that produced that change. 

To overwrite the current sha1sums, just run `python run.py > expected_output.txt`

