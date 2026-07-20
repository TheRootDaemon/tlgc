complete -c tlgc -s u -l update -d "Update the cache"
complete -c tlgc -s l -l list -d "List all pages in the current platform"
complete -c tlgc -s a -l list-all -d "List all pages"
complete -c tlgc -s s -l search -d "Search for pages containing a keyword"
complete -c tlgc -l list-platforms -d "List available platforms"
complete -c tlgc -l list-languages -d "List installed languages"
complete -c tlgc -s i -l info -d "Show cache information"
complete -c tlgc -s r -l render -d "Render the specified tldr page" -r
complete -c tlgc -l clean-cache -d "Interactively delete contents of the cache directory"
complete -c tlgc -l gen-config -d "Print the default config"
complete -c tlgc -l config-path -d "Print the default config path"
complete -c tlgc -s p -l platform -d "Specify the platform to use (linux, osx, windows, etc.)" -x -a \
    "(tlgc --offline --list-platforms 2> /dev/null)"
complete -c tlgc -s L -l language -d "Specify the languages to use" -x -a \
    "(tlgc --offline --list-languages 2> /dev/null)"
complete -c tlgc -l short-options -d "Display short options wherever possible (e.g. '-s')"
complete -c tlgc -l long-options -d "Display long options wherever possible (e.g. '--long')"
complete -c tlgc -l edit -d "Display a link to edit the shown page on GitHub"
complete -c tlgc -s o -l offline -d "Do not update the cache, even if it is stale"
complete -c tlgc -s c -l compact -d "Strip empty lines from output"
complete -c tlgc -l no-compact -d "Do not strip empty lines from output (overrides --compact)"
complete -c tlgc -s R -l raw -d "Print pages in raw markdown instead of rendering them"
complete -c tlgc -l no-raw -d "Render pages instead of printing raw file contents (overrides --raw)"
complete -c tlgc -s q -l quiet -d "Suppress status messages and warnings"
complete -c tlgc -l verbose -d "Be more verbose (can be specified twice)"
complete -c tlgc -l color -d "Specify when to enable color [default: auto] [possible values: auto, always, never]" -x -a "
    auto\t'Display color if standard output is a terminal and NO_COLOR is not set'
    always\t'Always display color'
    never\t'Never display color'
"
complete -c tlgc -l config -d "Specify an alternative path to the config file" -r
complete -c tlgc -s v -l version -d "Print version"
complete -c tlgc -s h -l help -d "Print help"
complete -c tlgc -f -a "(tlgc --offline --list-all 2> /dev/null)"
