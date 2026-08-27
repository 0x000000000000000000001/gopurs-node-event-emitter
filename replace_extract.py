import re
filepath = "/Users/0x1/Documents/htdocs/gopurs/gopurs-node-event-emitter/src/Node/EventEmitter.go"
with open(filepath, 'r') as f:
    content = f.read()

content = content.replace("ExtractEventEmitter", "nodeEventEmitter_extractEventEmitter")

with open(filepath, 'w') as f:
    f.write(content)
print("Replaced ExtractEventEmitter in EventEmitter.go")
