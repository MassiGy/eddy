#!/bin/bash

if [ ! -f  $HOME/.config/$(cat ./BINARY_NAME)-$(cat ./VERSION)/targets ]; then
    echo "Creating $HOME/.config/$(cat ./BINARY_NAME)-$(cat ./VERSION)/about file";

    figlet -W "         eddy" > $HOME/.config/$(cat ./BINARY_NAME)-$(cat ./VERSION)/about;
    echo "

        Eddy is a vim like editor, meaning that we have multiple motions 
        and modes. We currently support INSERT, NORMAL and PROMPT mode.

        INSERT is for updating the text buffer.
        NORMAL is mainly for navigation and editing using motions.
        PROMPT is for the command mode, it is like the ':' mode of vim.

    
        Some of eddy commands:  ( arrow keys are supported as well )

        Mode            Key         Name        Behavior
        ---------------------------------------------------------------
        NORMAL          i           insert      enter to INSERT mode.
        NORMAL          e           edit        enter to INSERT mode.
        NORMAL          ?           prompt      enter to PROMPT mode.
        NORMAL          :           prompt      enter to PROMPT mode.

        NORMAL          r           reload      reload the file. 
        NORMAL          w           write       save the file. 
        NORMAL          q           quit        quit the editor.

        NORMAL          h           vim.h       next character.
        NORMAL          l           vim.l       prev character.
        NORMAL          j           vim.j       goto next line.
        NORMAL          h           vim.h       goto prev line.
        NORMAL          f           vim.b       goto next word.
        NORMAL          b           vim.b       goto prev word.

        NORMAL          d           vim.db      delete prev word.
        NORMAL          D           vim.dw      delete next word.
        NORMAL          c           vim.cb      change prev word.
        NORMAL          C           vim.cw      change next word.

        NORMAL          I           vim.I       insert-in line start.
        NORMAL          A           vim.A       insert-in line end.
        NORMAL          o           vim.o       insert newline below.
        NORMAL          O           vim.O       insert newline above.
        NORMAL          y           vim.\"+yy    copy line to sysclip.
        NORMAL          p           vim.\"+p     paste line from sysclip.

        NORMAL          u           vim.u       undo (infinite)
        NORMAL          U           vim.CtrlR   redo (infinite)

        ALL             Ctrl+s      write       save the file. 
        ALL             Ctrl+r      reload      reload the file.
        ALL             Ctrl+q      quit        quit the editor
        ALL             Ctrl+n      vim.nu      toggle line numbers.
        ALL             Esc         escape      goto NORMAL mode.
        ALL             HOME        -           goto start of line.
        ALL             END         -           goto end of line.
        ALL             PgUp        page up     goto top of file.
        ALL             PgDown      page down   goto end of file.
    
    ">>$HOME/.config/$(cat ./BINARY_NAME)-$(cat ./VERSION)/about;
fi