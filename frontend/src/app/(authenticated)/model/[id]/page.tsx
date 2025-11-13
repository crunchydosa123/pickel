type PageProps = {
  params: {
    id: string;
  };
};

const Page = async ({ params }: PageProps) => {
  const { id } = await params;
  return (
    <div className="bg-blue-200 border flex flex-col p-2">
      <div className="text-3xl font-bold">Model Name: {id}</div>
      <div className="flex flex-col mt-10">
        <div>Deployments</div>
      </div>
    </div>
  )
};

export default Page;
