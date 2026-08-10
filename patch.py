import re
with open("src/Node/EventEmitter.go", "r") as f:
    text = f.read()

text = text.replace("gopurs_runtime.Apply(l.cb, arg1.(gopurs_runtime.Value))", "eff := gopurs_runtime.Apply(l.cb, arg1.(gopurs_runtime.Value))\n\t\tgo gopurs_runtime.Apply(eff, gopurs_runtime.Any(nil))")
text = text.replace("gopurs_runtime.UncurriedApp2(l.cb, arg1, arg2.(gopurs_runtime.Value))", "eff := gopurs_runtime.UncurriedApp2(l.cb, arg1, arg2.(gopurs_runtime.Value))\n\t\tgo gopurs_runtime.Apply(eff, gopurs_runtime.Any(nil))")
text = text.replace("gopurs_runtime.UncurriedApp3(l.cb, arg1, arg2, arg3.(gopurs_runtime.Value))", "eff := gopurs_runtime.UncurriedApp3(l.cb, arg1, arg2, arg3.(gopurs_runtime.Value))\n\t\tgo gopurs_runtime.Apply(eff, gopurs_runtime.Any(nil))")
text = text.replace("gopurs_runtime.UncurriedApp4(l.cb, arg1, arg2, arg3, arg4.(gopurs_runtime.Value))", "eff := gopurs_runtime.UncurriedApp4(l.cb, arg1, arg2, arg3, arg4.(gopurs_runtime.Value))\n\t\tgo gopurs_runtime.Apply(eff, gopurs_runtime.Any(nil))")
text = text.replace("gopurs_runtime.UncurriedApp5(l.cb, arg1, arg2, arg3, arg4, arg5.(gopurs_runtime.Value))", "eff := gopurs_runtime.UncurriedApp5(l.cb, arg1, arg2, arg3, arg4, arg5.(gopurs_runtime.Value))\n\t\tgo gopurs_runtime.Apply(eff, gopurs_runtime.Any(nil))")
text = text.replace("gopurs_runtime.UncurriedApp6(l.cb, arg1, arg2, arg3, arg4, arg5, arg6.(gopurs_runtime.Value))", "eff := gopurs_runtime.UncurriedApp6(l.cb, arg1, arg2, arg3, arg4, arg5, arg6.(gopurs_runtime.Value))\n\t\tgo gopurs_runtime.Apply(eff, gopurs_runtime.Any(nil))")

with open("src/Node/EventEmitter.go", "w") as f:
    f.write(text)
