# Contributing to fjira

🚀 Thanks for taking the time to look here! 🚀

In short - everyone is welcomed. Feel free to open PR with feature request or a bug fix.

### Branch name convention

`feature|fix|misc|docs/short-branch-name`

### Commit message convention

`feat|fix|misc|docs: my commit message goes here`

### Project structure

It's a classic golang project, with organized directories structure.
The most important directory is `internal` - it contains all the source code.
Why internal? because - at least for now - fjira is not intended to be a shared go module. It's just an app that makes
corporate rats (like me) life easier ;)

Internal structure looks like this:

```text
.
├── 📂 app  
│   └── ...
|   📂 fjira 
│   └── ...
|   📂 jira 
│   └── ...
```

#### app

It contains application engine, so everything that's needed in order to do "the thing" - whatever the thing is 😅.
You can notice that it doesn't contain any unit tests. There is a reason behind that. Just imagine that `app` module is
a vehicle. There is difference if you achieve your goal with ferrari, or old fiat. This is why it's not tested directly
here, but by another modules that contains business logic.

#### fjira

The heart of application - the main context - business logic of our application.

##### jira

Everything that's related to Jira, and Jira REST API.
