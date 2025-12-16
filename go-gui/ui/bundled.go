// Package main provides bundled resources (icons, images) for the application.
package ui

import (
	_ "embed"

	"fyne.io/fyne/v2"
)

//go:embed assets/icon.png
var iconPNG []byte

// IconResource is the application icon (PNG format required by Fyne).
var IconResource = &fyne.StaticResource{
	StaticName:    "icon.png",
	StaticContent: iconPNG,
}

// LogoResource is the full logo for display in the UI header (SVG format).
var LogoResource = &fyne.StaticResource{
	StaticName: "logo.svg",
	StaticContent: []byte(`<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 512 512">
  <defs>
    <linearGradient id="chipGrad" x1="0%" y1="0%" x2="100%" y2="100%">
      <stop offset="0%" style="stop-color:#1C1C26"/>
      <stop offset="100%" style="stop-color:#12121A"/>
    </linearGradient>
    <filter id="glow" x="-50%" y="-50%" width="200%" height="200%">
      <feGaussianBlur stdDeviation="4" result="coloredBlur"/>
      <feMerge>
        <feMergeNode in="coloredBlur"/>
        <feMergeNode in="SourceGraphic"/>
      </feMerge>
    </filter>
    <filter id="glowStrong" x="-50%" y="-50%" width="200%" height="200%">
      <feGaussianBlur stdDeviation="8" result="coloredBlur"/>
      <feMerge>
        <feMergeNode in="coloredBlur"/>
        <feMergeNode in="SourceGraphic"/>
      </feMerge>
    </filter>
  </defs>
  <circle cx="256" cy="256" r="240" fill="#12121A" stroke="#3C3C50" stroke-width="4"/>
  <circle cx="256" cy="256" r="220" fill="none" stroke="#00D4FF" stroke-width="2" opacity="0.3"/>
  <rect x="120" y="120" width="272" height="272" rx="20" ry="20" fill="url(#chipGrad)" stroke="#3C3C50" stroke-width="3"/>
  <g fill="#00D4FF" filter="url(#glow)">
    <rect x="160" y="100" width="12" height="30" rx="2"/><rect x="190" y="100" width="12" height="30" rx="2"/>
    <rect x="220" y="100" width="12" height="30" rx="2"/><rect x="280" y="100" width="12" height="30" rx="2"/>
    <rect x="310" y="100" width="12" height="30" rx="2"/><rect x="340" y="100" width="12" height="30" rx="2"/>
    <rect x="160" y="382" width="12" height="30" rx="2"/><rect x="190" y="382" width="12" height="30" rx="2"/>
    <rect x="220" y="382" width="12" height="30" rx="2"/><rect x="280" y="382" width="12" height="30" rx="2"/>
    <rect x="310" y="382" width="12" height="30" rx="2"/><rect x="340" y="382" width="12" height="30" rx="2"/>
    <rect x="100" y="160" width="30" height="12" rx="2"/><rect x="100" y="190" width="30" height="12" rx="2"/>
    <rect x="100" y="220" width="30" height="12" rx="2"/><rect x="100" y="280" width="30" height="12" rx="2"/>
    <rect x="100" y="310" width="30" height="12" rx="2"/><rect x="100" y="340" width="30" height="12" rx="2"/>
    <rect x="382" y="160" width="30" height="12" rx="2"/><rect x="382" y="190" width="30" height="12" rx="2"/>
    <rect x="382" y="220" width="30" height="12" rx="2"/><rect x="382" y="280" width="30" height="12" rx="2"/>
    <rect x="382" y="310" width="30" height="12" rx="2"/><rect x="382" y="340" width="30" height="12" rx="2"/>
  </g>
  <rect x="150" y="150" width="212" height="212" rx="10" fill="#1C1C26" stroke="#3C3C50" stroke-width="2"/>
  <g stroke="#00D4FF" stroke-width="3" fill="none" filter="url(#glowStrong)">
    <path d="M 170 230 h 20 v -20 h 20 v 20 h 20 v -20 h 20 v 20 h 20"/>
    <path d="M 342 280 h -20 v 20 h -20 v -20 h -20 v 20 h -20 v -20 h -20"/>
  </g>
  <text x="175" y="200" fill="#00FF88" font-family="monospace" font-size="18" font-weight="bold" filter="url(#glow)">TX</text>
  <text x="315" y="320" fill="#FF5252" font-family="monospace" font-size="18" font-weight="bold" filter="url(#glow)">RX</text>
  <text x="256" y="265" fill="#00D4FF" font-family="Arial, sans-serif" font-size="42" font-weight="bold" text-anchor="middle" filter="url(#glowStrong)">PS3</text>
  <text x="256" y="310" fill="#8C8CA0" font-family="Arial, sans-serif" font-size="28" font-weight="bold" text-anchor="middle">SC</text>
  <g stroke="#00D4FF" stroke-width="2" fill="none" opacity="0.5">
    <path d="M 80 140 L 80 80 L 140 80"/><circle cx="80" cy="80" r="4" fill="#00D4FF"/>
    <path d="M 432 140 L 432 80 L 372 80"/><circle cx="432" cy="80" r="4" fill="#00D4FF"/>
    <path d="M 80 372 L 80 432 L 140 432"/><circle cx="80" cy="432" r="4" fill="#00D4FF"/>
    <path d="M 432 372 L 432 432 L 372 432"/><circle cx="432" cy="432" r="4" fill="#00D4FF"/>
  </g>
  <circle cx="165" cy="365" r="6" fill="#00FF88" filter="url(#glow)"/>
  <circle cx="185" cy="365" r="6" fill="#00D4FF" filter="url(#glow)"/>
  <circle cx="327" cy="365" r="6" fill="#FFAA00" filter="url(#glow)" opacity="0.6"/>
  <circle cx="347" cy="365" r="6" fill="#FF5252" filter="url(#glow)" opacity="0.4"/>
</svg>`),
}
