"use client";

import { useState } from "react";
import { Button } from "@/components/ui/button";
import { Select, SelectTrigger, SelectContent, SelectItem, SelectValue } from "@/components/ui/select";

type Repo = {
  full_name: string;
};

type Props = {
  model: any;
  repos: Repo[];
};

export default function RepoDropdownClient({ model, repos }: Props) {
  const [selectedRepo, setSelectedRepo] = useState("");

  const openInstallWindow = () => {
    const INSTALL_URL = "https://github.com/apps/pickel-deploy-bot/installations/new";
    window.open(
      INSTALL_URL,
      "githubInstall",
      "width=800,height=700,menubar=no,toolbar=no,location=no,status=no"
    );
  };

  const linkRepoToModel = async () => {
    if (!selectedRepo) {
      alert("Select a repo first!");
      return;
    }

    try {
      await fetch("/api/link-repo", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({
          model_id: model.id,
          repo_full_name: selectedRepo,
        }),
      });
      alert("Repo linked!");
    } catch (err) {
      console.error(err);
      alert("Failed to link repo");
    }
  };

  return (
    <div className="flex flex-col space-y-2 w-80">
      <Button onClick={openInstallWindow}>Connect GitHub</Button>
      <Select onValueChange={setSelectedRepo}>
        <SelectTrigger>
          <SelectValue placeholder="Select repo" />
        </SelectTrigger>
        <SelectContent>
          {repos.map((repo) => (
            <SelectItem key={repo.full_name} value={repo.full_name}>
              {repo.full_name}
            </SelectItem>
          ))}
        </SelectContent>
      </Select>
      <Button onClick={linkRepoToModel}>Link Repo</Button>
    </div>
  );
}
