import { Button } from "@/components/ui/button";
import {
  CalendarIcon,
  HomeIcon,
  LibraryIcon,
  ListIcon,
  LogInIcon,
  StarIcon,
} from "./icons";

export default function NavigationBar() {
  return (
    <aside className="flex w-64 flex-col bg-gray-100 p-4">
      <nav className="flex flex-col space-y-4">
        <Button variant="ghost" className="flex items-center space-x-2">
          <HomeIcon className="h-6 w-6" />
          <span>Home</span>
        </Button>
        <Button variant="ghost" className="flex items-center space-x-2">
          <StarIcon className="h-6 w-6" />
          <span>Favorites</span>
        </Button>
        <Button variant="ghost" className="flex items-center space-x-2">
          <LibraryIcon className="h-6 w-6" />
          <span>Library</span>
        </Button>
        <Button variant="ghost" className="flex items-center space-x-2">
          <CalendarIcon className="h-6 w-6" />
          <span>Calendar</span>
        </Button>
        <Button variant="ghost" className="flex items-center space-x-2">
          <ListIcon className="h-6 w-6" />
          <span>Project List</span>
        </Button>
      </nav>
      <Button variant="ghost" className="mt-auto flex items-center space-x-2">
        <LogInIcon className="h-6 w-6" />
        <span>Login</span>
      </Button>
    </aside>
  );
}
