export default function Home() {
  return (
    <main className="p-8">
      <h1 className="text-4xl font-bold mb-4">Welcome to SpecScope</h1>
      <p className="mb-6">Explore RF spectrum activity, detect interference, and more.</p>
      <div className="space-x-4">
        <a href="/dashboard" className="bg-blue-600 text-white px-4 py-2 rounded">Go to Dashboard</a>
        <a href="/upload" className="bg-gray-700 text-white px-4 py-2 rounded">Upload Data</a>
      </div>
    </main>
  );
}
