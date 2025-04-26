import { Button } from "@/components/ui/button"; // Import Shadcn Button

function App() {
  return (
    <div className="p-4"> {/* Added Tailwind padding class */}
      <h1 className="text-2xl font-bold text-blue-600"> {/* Added Tailwind text classes */}
        GoSocial Frontend
      </h1>
      <p>Tailwind CSS is working!</p>
      <Button className="mt-4">Shadcn Button</Button> {/* Added Shadcn Button */}
    </div>
  )
}

export default App
