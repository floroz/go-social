import { zodResolver } from "@hookform/resolvers/zod";
import { useForm } from "react-hook-form";
import { z } from "zod";
import { useNavigate } from "react-router";

import { Button } from "@/components/ui/button";
import {
  Form,
  FormControl,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from "@/components/ui/form";
import { Input } from "@/components/ui/input";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { LoginRequestDTO } from "@/domain/auth.dto";
import { useLogin } from "@/hooks/useLogin"; 

// Define Zod schema for login
const formSchema = z.object({
  email: z.string().email({ message: "Invalid email address." }),
  password: z.string().min(3, { message: "Password is required." }), // we let the backend handle length 
});

const LoginPage = () => {
  const navigate = useNavigate();

  // Use the custom hook for login logic
  const { login, isLoading, error } = useLogin({
    onSuccess: (user) => {
      // Side effect specific to this page: navigate after successful login
      console.log("Login successful on page:", user);
      navigate('/'); 
      // TODO: Show success toast
    },
    onError: (error) => {
      console.error("Login failed on page:", error.message);
      // TODO: Show error toast
    },
  });

  const form = useForm<z.infer<typeof formSchema>>({
    resolver: zodResolver(formSchema),
    defaultValues: {
      email: "",
      password: "",
    },
  });

  function onSubmit(values: z.infer<typeof formSchema>) {
    const loginData: LoginRequestDTO = {
      email: values.email,
      password: values.password,
    };
    login(loginData); 
  }

  return (
    <div className="container mx-auto flex h-screen items-center justify-center p-4">
      <Card className="w-full max-w-sm"> 
        <CardHeader>
          <CardTitle>Login</CardTitle>
          <CardDescription>Enter your credentials to access your account.</CardDescription>
          {error && <p className="text-sm font-medium text-destructive">{error.message}</p>}
        </CardHeader>
        <CardContent>
          <Form {...form}>
            <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-4">
              <FormField
                control={form.control}
                name="email"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>Email</FormLabel>
                    <FormControl>
                      <Input type="email" placeholder="you@example.com" {...field} />
                    </FormControl>
                    <FormMessage />
                  </FormItem>
                )}
              />
              <FormField
                control={form.control}
                name="password"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>Password</FormLabel>
                    <FormControl>
                      <Input type="password" placeholder="********" {...field} />
                    </FormControl>
                    <FormMessage />
                  </FormItem>
                )}
              />
              <Button type="submit" className="w-full" disabled={isLoading}>
                {isLoading ? 'Logging In...' : 'Login'}
              </Button>
            </form>
          </Form>
        </CardContent>
      </Card>
    </div>
  );
};

export default LoginPage;
