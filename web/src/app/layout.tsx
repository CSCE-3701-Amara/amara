import type { Metadata } from "next";
import localFont from "next/font/local";
import "./globals.css";

const funnelSans = localFont({
  src: "../../public/fonts/FunnelSans.ttf",
  variable: "--font-funnel-sans",
  weight: "300 800",
  display: "swap",
});

const funnelDisplay = localFont({
  src: "../../public/fonts/FunnelDisplay.ttf",
  variable: "--font-funnel-display",
  weight: "300 800",
  display: "swap",
});

export const metadata: Metadata = {
  title: {
    default: "amara — Skills-first recruitment",
    template: "%s | amara",
  },
  description:
    "Find opportunities that value your skills. Build better hiring teams with amara.",
};

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <html
      lang="en"
      className={`${funnelSans.variable} ${funnelDisplay.variable} h-full antialiased`}
    >
      <body className="min-h-full flex flex-col">{children}</body>
    </html>
  );
}
