import * as htmx from 'htmx.org';

import { createIcons, Mail, Lock, Repeat, LogOut, DiamondPercent, ChevronDown, Bell, Gauge } from 'lucide';

const LoadIcons = () => createIcons({
  icons: {
    Mail,
    Lock,
    Repeat,
    LogOut,
    DiamondPercent,
    ChevronDown,
    Bell,
    Gauge
  }
})
window.htmx = htmx;

import preload from "./htmx-preload"
document.addEventListener("DOMContentLoaded", function (event) {
  LoadIcons();
});

document.addEventListener('htmx:load', function (evt) {
  preload()
  LoadIcons();
})
