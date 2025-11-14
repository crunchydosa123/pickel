import AddModelCode from "@/components/AddModelCode";
import { Button } from "@/components/ui/button";
import { cookies } from "next/headers";

type PageProps = {
  params: {
    id: string;
  };
};

const Page = async ({ params }: PageProps) => {
  const { id } = await params;

  const baseUrl = process.env.NEXT_PUBLIC_BASE_URL || "http://localhost:3000";

  const cookieStore = await cookies();
    const cookieHeader = cookieStore
      .getAll()
      .map((c) => `${c.name}=${c.value}`)
      .join('; ');
  
    let data: any = {};

  try {
    const res = await fetch(`${baseUrl}/api/model/${id}`,{
      method: "GET",
      headers: {
        Cookie: cookieHeader,
      },
      cache: 'no-store',
    }
    )
    const data = await res.json();
console.log("DATA:", data);
  }catch(err){
    console.error(err);
  }

  return (
    <div className="bg-blue-200 border flex flex-col p-2">
      <div className="text-3xl font-bold">Model Name: {id}</div>
      <div className="flex flex-col mt-10">
        <div className="flex justify-between">
          <div>Deployments</div>
          <AddModelCode />
        </div>
        
      </div>
    </div>
  )
};

export default Page;
