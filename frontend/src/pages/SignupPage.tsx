import { zodResolver } from "@hookform/resolvers/zod";
import { useForm } from "react-hook-form";
import { z } from "zod";
import { useNavigate, Link } from "react-router";

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
import { SignupRequest } from "@/types/api";
import { useSignup } from "@/hooks/useSignup";

// Define Zod schema based on backend validation rules
// (Schema remains the same)
const formSchema = z.object({
  firstName: z.string().min(3, { message: "First name must be at least 3 characters." }).max(50),
  lastName: z.string().min(3, { message: "Last name must be at least 3 characters." }).max(50),
  username: z.string().min(3).max(50).regex(/^[a-zA-Z0-9]+$/, { message: "Username must be alphanumeric." }),
  email: z.string().email({ message: "Invalid email address." }).min(3).max(50),
  password: z.string().min(8, { message: "Password must be at least 8 characters." }).max(50),
});

const SignupPage = () => {
  const navigate = useNavigate(); 

  const { signup, isLoading, error } = useSignup({
    onSuccess: (user) => {
      navigate('/'); 
      console.log("Signup successful on page:", user);
      // TODO: Show success toast
    },
    onError: (error) => {
      console.error("Signup failed on page:", error.message);
      // TODO: Show error toast
    },
  });

  const form = useForm<z.infer<typeof formSchema>>({
    resolver: zodResolver(formSchema),
    defaultValues: {
      firstName: "",
      lastName: "",
      username: "",
      email: "",
      password: "",
    },
  });

  function onSubmit(values: z.infer<typeof formSchema>) {
    const backendData: SignupRequest = {
      first_name: values.firstName,
      last_name: values.lastName,
      username: values.username,
      email: values.email,
      password: values.password,
    };
    signup(backendData); 
  }

  return (
    <div className="container mx-auto flex h-screen items-center justify-center p-4">
      <Card className="w-full max-w-md">
        <CardHeader>
          <CardTitle>Sign Up</CardTitle>
          <CardDescription>Create your GoSocial account.</CardDescription>
          {error && <p className="text-sm font-medium text-destructive">{error.message}</p>}
        </CardHeader>
        <CardContent>
          <Form {...form}>
            <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-4">
              <FormField
                control={form.control}
                name="firstName"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>First Name</FormLabel>
                    <FormControl>
                      <Input placeholder="John" {...field} />
                    </FormControl>
                    <FormMessage />
                  </FormItem>
                )}
              />
              <FormField
                control={form.control}
                name="lastName"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>Last Name</FormLabel>
                    <FormControl>
                      <Input placeholder="Doe" {...field} />
                    </FormControl>
                    <FormMessage />
                  </FormItem>
                )}
              />
              <FormField
                control={form.control}
                name="username"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>Username</FormLabel>
                    <FormControl>
                      <Input placeholder="johndoe" {...field} />
                    </FormControl>
                    <FormMessage />
                  </FormItem>
                )}
              />
              <FormField
                control={form.control}
                name="email"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>Email</FormLabel>
                    <FormControl>
                      <Input type="email" placeholder="john.doe@example.com" {...field} />
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
                {isLoading ? 'Signing Up...' : 'Sign Up'}
              </Button>
            </form>
          </Form>
          <p className="mt-4 text-center text-sm text-muted-foreground">
            Already have an account?{' '}
            <Link to="/login" className="underline underline-offset-4 hover:text-primary">
              Log in
            </Link>
          </p>
        </CardContent>
      </Card>
    </div>
  );
};

export default SignupPage;
