import pathlib

# All Chinese text replacements organized by file
fixes = {
    "src/components/Heatmap.vue": [
        (b'= "?', b'= "'),
        ('" + (cm.value + 1) + "?"', '" + (cm.value + 1) + "\u5e74"'),
        ('" + count + " ??"', '" + count + " \u6761\u5907\u5fd8\u5f55"'),
    ],
    "src/components/NoteCard.vue": [
        ('"??"', '"\u5929"'),
        ('"???"', '"\u5c0f\u65f6"'),
        # These need careful matching
    ],
    "src/components/LoginDialog.vue": [
        ('"???', '"\u7528\u6237\u540d'),
        ('???"', '\u5bc6\u7801"'),
    ],
}

for fname, replacements in fixes.items():
    p = pathlib.Path("F:/vue/Meng") / fname
    c = p.read_text("utf-8")
    for old, new in replacements:
        if isinstance(old, bytes):
            c_bytes = p.read_bytes()
            if isinstance(new, bytes):
                c_bytes = c_bytes.replace(old, new)
            else:
                c_bytes = c_bytes.replace(old, new.encode("utf-8"))
            p.write_bytes(c_bytes)
        else:
            c = c.replace(old, new)
            p.write_text(c, "utf-8")
    print(f"Fixed {fname}")

print("Done")
