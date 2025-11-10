import type { Metadata } from "next";
import "./globals.css";

export const metadata: Metadata = {
  title: "Scrypts - Secure Notes",
  description: "Your thoughts, encrypted.",
};

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <html lang="en" className="dark">
      <body className="font-mono bg-background text-foreground antialiased">
        {children}
      </body>
    </html>
  );
}
