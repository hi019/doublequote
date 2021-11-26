import { Nav as NavElement, NavItem } from "../components/Nav";
import React from "react";
import { useAppSelector } from "../../store/hooks";

const navItems: NavItem[] = [
  // nav items
  {
    name: "Collections",
    href: "#",
    current: false,
    profile: false,
    requiresSignIn: true,
  },

  // Profile items
  {
    name: "Settings",
    href: "#",
    current: false,
    profile: true,
  },
  {
    name: "Sign out",
    href: "#",
    current: false,
    profile: true,
  },
];

export const Nav = () => {
  // TODO selecting menu items
  const [items] = React.useState(navItems);
  const isSignedIn = useAppSelector((s) => s.user.isSignedIn);

  return <NavElement items={items} isSignedIn={isSignedIn ?? false} />;
};
