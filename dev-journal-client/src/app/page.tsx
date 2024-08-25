import { Button } from "@/components/ui/button";
import { PlusIcon } from "./_components/icons";

export default function HomePage() {
  return (
    <main className="flex-1 p-4">
      <div className="mb-4 flex items-center justify-between">
        <h2 className="text-xl font-semibold">Recent Notes</h2>
        <Button variant="outline" className="flex items-center space-x-2">
          <PlusIcon className="h-6 w-6" />
          <span>Add New</span>
        </Button>
      </div>
      <div className="rounded-lg bg-gray-50 p-4 text-center">
        <p className="text-lg font-semibold">
          Nothing here yet. Add a project to get started!
        </p>
        <p className="text-gray-500">
          Your most recent notes will show up here. Click on &apos;Add new&apos;
          to get started!
        </p>
      </div>
    </main>
  );
}
