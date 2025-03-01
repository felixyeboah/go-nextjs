"use client";

import { useState } from "react";
import { useRouter } from "next/navigation";
import { zodResolver } from "@hookform/resolvers/zod";
import { useForm } from "react-hook-form";
import * as z from "zod";
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
import { useToast } from "@/components/ui/use-toast";

const resendVerificationSchema = z.object({
  email: z.string().email("Please enter a valid email address"),
});

type ResendVerificationFormValues = z.infer<typeof resendVerificationSchema>;

export function ResendVerificationForm() {
  const router = useRouter();
  const { toast } = useToast();
  const [isLoading, setIsLoading] = useState(false);
  const [emailSent, setEmailSent] = useState(false);

  const form = useForm<ResendVerificationFormValues>({
    resolver: zodResolver(resendVerificationSchema),
    defaultValues: {
      email: "",
    },
  });

  async function onSubmit(data: ResendVerificationFormValues) {
    setIsLoading(true);

    try {
      const response = await fetch("/api/auth/resend-verification", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify(data),
      });

      if (!response.ok) {
        const errorData = await response.json();
        throw new Error(
          errorData.message || "Failed to resend verification email"
        );
      }

      setEmailSent(true);
      toast({
        title: "Verification email sent",
        description: "Please check your email for the verification link",
      });
    } catch (error) {
      toast({
        title: "Error",
        description:
          error instanceof Error
            ? error.message
            : "Failed to resend verification email",
        variant: "destructive",
      });
    } finally {
      setIsLoading(false);
    }
  }

  if (emailSent) {
    return (
      <div className="flex flex-col space-y-4">
        <div className="text-center">
          <p className="text-sm text-muted-foreground">
            Verification email sent! Please check your inbox.
          </p>
        </div>
        <Button onClick={() => setEmailSent(false)}>Try another email</Button>
      </div>
    );
  }

  return (
    <Form {...form}>
      <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-6">
        <FormField
          control={form.control}
          name="email"
          render={({ field }) => (
            <FormItem>
              <FormLabel>Email</FormLabel>
              <FormControl>
                <Input
                  placeholder="name@example.com"
                  type="email"
                  autoCapitalize="none"
                  autoComplete="email"
                  autoCorrect="off"
                  disabled={isLoading}
                  {...field}
                />
              </FormControl>
              <FormMessage />
            </FormItem>
          )}
        />
        <Button type="submit" className="w-full" disabled={isLoading}>
          {isLoading ? "Sending..." : "Send verification email"}
        </Button>
      </form>
    </Form>
  );
}
