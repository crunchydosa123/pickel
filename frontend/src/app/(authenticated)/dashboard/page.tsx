import { cookies } from 'next/headers';
import CreateModelPopover from '@/components/CreateModelPopover';
import { Card } from '@/components/ui/card';

export default async function Page() {
  const baseUrl = process.env.NEXT_PUBLIC_BASE_URL || 'http://localhost:3000';

  // Convert Next.js cookies() to a valid Cookie header string
  const cookieStore = await cookies();
  const cookieHeader = cookieStore
    .getAll()
    .map((c) => `${c.name}=${c.value}`)
    .join('; ');

  let data: any = {};

  try {
    const res = await fetch(`${baseUrl}/api/model/getAll`, {
      method: 'GET',
      headers: {
        Cookie: cookieHeader,
      },
      cache: 'no-store',
    });

    data = await res.json();
  } catch (err) {
    console.error('Failed to fetch models:', err);
    data = { error: 'Failed to fetch models' };
  }

  const models = data.models || [];

  return (
    <div className="w-full min-h-screen p-6">
      <div className="flex justify-between items-center mb-4">
        <h1 className="text-2xl font-bold text-black">Your Models</h1>
        <CreateModelPopover />
      </div>

      {models.length > 0 ? (
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
          {models.map((model: any) => (
            <Card
              key={model.id}
              className="bg-white text-gray-800 rounded-lg p-4 transition-shadow"
            >
              <h2 className="font-semibold text-lg">{model.name}</h2>
              <p className="text-sm text-gray-600 mt-1">{model.id}</p>
            </Card>
          ))}
        </div>
      ) : (
        <p className="text-gray-100 mt-4">
          No models found. Create one using the button above.
        </p>
      )}
    </div>
  );
}
