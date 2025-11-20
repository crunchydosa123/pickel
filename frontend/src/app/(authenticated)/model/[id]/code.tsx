"use client"

import { Button } from "@/components/ui/button";
import { Github } from "lucide-react";

type Props = {
  model: any;
};

const ModelCode = ({ model }: Props) => {
  const INSTALL_URL = "https://github.com/apps/pickel-deploy-bot/installations/new";

  return (
    <div className="w-full bg-gray-300 p-4 h-full flex items-center">
      {model.deployment_type === "0" && (
        <div className="flex w-full items-center justify-center">
          <Button
            className="gap-2 px-6 py-4 text-xl font-semibold"
            onClick={() => window.location.href = INSTALL_URL}
          >
            <Github className="w-6 h-6" />
            Connect GitHub
          </Button>
        </div>
      )}
    </div>
  );
};

export default ModelCode;
