def _reset_to_host_impl(settings, attr):
    # Setting an empty list resets the target platform requested.
    # Thus, the host platform will be used.
    return {"//command_line_option:platforms": []}

reset_to_host_transition = transition(
    implementation = _reset_to_host_impl,
    inputs = [],
    outputs = ["//command_line_option:platforms"],
)

def _host_configured_target_impl(ctx):
    # Just pass through the files from the underlying target
    return [DefaultInfo(files = depset(ctx.files.target))]

host_configured_target = rule(
    implementation = _host_configured_target_impl,
    attrs = {
        "target": attr.label(cfg = reset_to_host_transition),
    },
)
