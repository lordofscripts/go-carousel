# Go Carousel For Developers


This document is meant as a general guide to extend the Go Carousel application,
as well as for maintenance purposes. You never know if after a long time doing
some other work, you would remember old code, right? At least that is my 
philosophy.

## Session Handlers

Although originally `goCarousel` was a Bash script that worked with the Gnome
session manager, it now supports other session managers. After several months
of use I decided to port the Bash script to a pure GO application. After all
shell scripts get squeaky the longer they get. And I prefer real programming
languages too.

The original `gnomeChangeBackground` then renamed `goCarousel` supported only
the **Gnome** session manager. But today I decided to support XFCE as well. As
a result, I wanted a way to let it autodetermine which method to use, and then
which external application and parameters to use to accomplish the task. As you
know, modern session managers use communication channels to work with their
configuration rather than directly changing configuration files.

*Go Carousel* uses the `GDMSESSION` environment variable to determine which
session handler to use (Gnome, XFCE). There are built-in handlers for each of
the supported session managers.

### Gnome Sessions

Tested with Gnome 43 & Gnome 48; therefore, I assume it should work with all
Gnome versions in between. Since *Cinnamon*  also uses `gsettings` this handler
works with both.

**Gnome** knows a Dark & Light color theme and uses the `gsettings` application
to tweak settings.

`gsettings get org.gnome.desktop.interface color-scheme` is the command used
to query the current *color scheme*. For example, I get `prefer-dark`but it
depends on what the user has set.

`gsettings set org.gnome.desktop.background picture-uri file://path/to/wallpaper.png` 
is used to set the *Light* color scheme wallpaper.

`gsettings set org.gnome.desktop.background picture-uri-dark file://path/to/wallpaper.png` 
is used to set the *Dark* color scheme wallpaper.

#### Cinnamon Sessions

Uses the Gnome Sessions handler with a different schemas `org.cinnamon.desktop.interface`
and `org.cinnamon.desktop.background`. And the `color-scheme` key is replaced by the
`gtk-theme` key.

### XFCE Sessions

Tested with XFCE v4. This session manager uses the `xfce4-desktop` communication channel
via the `xfconf-query` tool. When using a GUI XFCE4 can be tweaked with the
`xfce4-settings-editor` application. It shows its communication channels and their
respective properties.

`xfconf-query -c xfce4-desktop -l` lists all properties of the `xfce4-desktop` channel. As
we can see, it is aware of multiple screens and multiple monitors for every screen.
Note: *Go Carousel only works with screen0/monitor0 properties*.

`xfconf-query -c xsettings -p /Net/ThemeName` produces the name of the current theme. For 
example it can return `Xfce` in the default configuration. As for the theme variants, some
themes have variants. For example `Adwaita-dark` (dark) and `Adwaita` (light). The
built-in handler uses the presence of the word "dark" to detect the *Dark* variant.

`xfconf-query -c xfce4-desktop -p /backdrop/screen0/monitor0/color-style` shows us the
current color scheme. In my case it returns `1`

