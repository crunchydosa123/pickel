import React from 'react'
import { cookies } from 'next/headers';
import RepoDropdown from './repo-dropdown';

type Props = {}

type Repo = {
  full_name: string;
};


export async function Repos (props: Props) {
  const cookieStore = await cookies();
  const cookieHeader = cookieStore
    .getAll()
    .map((c) => `${c.name}=${c.value}`)
    .join("; ");

  let repos: Repo[] = [];

  // Fetch installed repos for the user
  try {
    const res = await fetch("http://localhost:8080/api/github/installed-repos", {
      method: "GET",
      headers: {
        Cookie: cookieHeader,
      },
      cache: "no-store",
    });
    
    const data = await res.json();
    repos = data || [];
    console.log("Fetched repos:", res); // Server-side log
  } catch (err) {
    console.error("Failed to fetch installed repos:", err);
  }
  return (
   <RepoDropdown repos={repos} onLinkRepo={() => {}} />
  )
}

export default Repos