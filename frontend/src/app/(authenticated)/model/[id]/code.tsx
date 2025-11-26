import RepoDropdown from './repo-dropdown';

type Props = {
  model: any;
};

type Repo = {
  full_name: string;
};

export default async function ModelCodeServer({ model }: Props) {
  // Server-side cookie access
  const baseUrl = process.env.NEXT_PUBLIC_BASE_URL || 'http://localhost:3000';
  const { cookies } = await import('next/headers');
  const cookieStore = await cookies();
  const cookieHeader = cookieStore
    .getAll()
    .map((c) => `${c.name}=${c.value}`)
    .join("; ");

  // Fetch installed repos
  let repos: Repo[] = [];

  try {
    const res = await fetch(`${baseUrl}/api/github/installed-repos`, {
      method: "GET",
      headers: { Cookie: cookieHeader },
      cache: "no-store",
    });

    // Get raw text first for debugging
    const text = await res.text();
    try {
      const data = JSON.parse(text);
      repos = (data.repositories || []) as Repo[];
      console.log("Fetched repos on server:", repos);
    } catch (jsonErr) {
      console.error("Failed to parse JSON:", jsonErr, "Response text:", text);
    }
  } catch (err) {
    console.error("Fetch failed:", err);
  }

  return (
    <div>
      <h2>Model: {model.id}</h2>
      <RepoDropdown model={model} repos={repos} />
    </div>
  );
}
