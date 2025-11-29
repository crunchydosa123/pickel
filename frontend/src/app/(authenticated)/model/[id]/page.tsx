import AddModelCode from "@/components/AddModelCode";
import { Button } from "@/components/ui/button";
import { cookies } from "next/headers";
import ModelTabs from "./model-tabs";

type Props = {
  params: Promise<{ id: string }>;
};

const Page = async ({ params }: Props) => {
  const { id } = await params; // must await

  const baseUrl = process.env.NEXT_PUBLIC_BASE_URL || "http://localhost:3000";

  const cookieStore = await cookies();
  const cookieHeader = cookieStore
    .getAll()
    .map((c) => `${c.name}=${c.value}`)
    .join('; ');

  let data: any = {};
  try {
    const res = await fetch(`${baseUrl}/api/model/${id}`, {
      method: "GET",
      headers: { Cookie: cookieHeader },
      cache: "no-store",
    });
    data = await res.json();
  } catch (err) {
    console.error(err);
  }

  return (
    <div className="bg flex flex-col p-5 justify-center">
      <div className="text-3xl font-bold">Model Name: {data?.name}</div>
      <div className="flex flex-col mt-10 w-full h-screen">
        <ModelTabs model={data} />
      </div>
    </div>
  );
};

export default Page;
