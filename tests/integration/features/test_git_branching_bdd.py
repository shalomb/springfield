import subprocess
import pytest
import os
import shutil
import tempfile
from pytest_bdd import scenario, given, when, then, parsers

@scenario('../../../features/git_branching.feature', 'Create a Feature Branch')
def test_create_feature_branch():
    pass

@scenario('../../../features/git_branching.feature', 'Create a Fix Branch')
def test_create_fix_branch():
    pass

@scenario('../../../features/git_branching.feature', 'Enforce Naming Convention')
def test_enforce_naming_convention():
    pass

@pytest.fixture
def test_repo():
    # Setup a temp git repo
    temp_dir = tempfile.mkdtemp()
    orig_dir = os.getcwd()
    try:
        os.chdir(temp_dir)
        subprocess.run(["git", "init"], check=True)
        # Configure git user for the temp repo
        subprocess.run(["git", "config", "user.email", "herb@example.com"], check=True)
        subprocess.run(["git", "config", "user.name", "Herb Powell"], check=True)
        
        with open("README.md", "w") as f:
            f.write("test")
        subprocess.run(["git", "add", "README.md"], check=True)
        subprocess.run(["git", "commit", "-m", "initial commit"], check=True)
        # Create a Justfile in the temp repo
        shutil.copy(os.path.join(orig_dir, "Justfile"), "Justfile")
        subprocess.run(["git", "add", "Justfile"], check=True)
        subprocess.run(["git", "commit", "-m", "add Justfile"], check=True)
        yield temp_dir
    finally:
        os.chdir(orig_dir)
        shutil.rmtree(temp_dir)

@pytest.fixture
def context():
    return {}

@given('I am on the "main" branch')
def on_main(test_repo):
    # Ensure we are on main (git init usually creates master or main)
    # Check current branch
    result = subprocess.run(["git", "branch", "--show-current"], capture_output=True, text=True)
    current = result.stdout.strip()
    if current != "main":
        subprocess.run(["git", "checkout", "-b", "main"], check=True)

@given('the working directory is clean')
def directory_clean():
    result = subprocess.run(["git", "status", "--porcelain"], capture_output=True, text=True)
    assert result.stdout.strip() == ""

@when(parsers.parse('I run "just start-feature \'{name}\'"'))
def run_start_feature(context, name):
    result = subprocess.run(["just", "start-feature", name], capture_output=True, text=True)
    context["result"] = result

@then(parsers.parse('a new branch "feat/{name}" should be created'))
def check_branch_created(name):
    result = subprocess.run(["git", "branch"], capture_output=True, text=True)
    assert f"feat/{name}" in result.stdout

@then(parsers.parse('I should be switched to "feat/{name}"'))
def check_switched(name):
    result = subprocess.run(["git", "branch", "--show-current"], capture_output=True, text=True)
    assert result.stdout.strip() == f"feat/{name}"

@when(parsers.parse('I run "just start-fix \'{name}\'"'))
def run_start_fix(context, name):
    result = subprocess.run(["just", "start-fix", name], capture_output=True, text=True)
    context["result"] = result

@then(parsers.parse('a new branch "fix/{name}" should be created'))
def check_fix_branch_created(name):
    result = subprocess.run(["git", "branch"], capture_output=True, text=True)
    assert f"fix/{name}" in result.stdout

@then(parsers.parse('I should be switched to "fix/{name}"'))
def check_fix_switched(name):
    result = subprocess.run(["git", "branch", "--show-current"], capture_output=True, text=True)
    assert result.stdout.strip() == f"fix/{name}"

@then('the command should fail')
def check_fail(context):
    assert context["result"].returncode != 0

@then(parsers.parse('the error message should mention "{text}"'))
def check_error_message(context, text):
    assert text in context["result"].stdout or text in context["result"].stderr
