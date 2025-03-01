"use client";

import { useEffect, useState } from "react";
import { useToast } from "@/components/ui/use-toast";
import {
  Table,
  TableBody,
  TableCaption,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";
import { Badge } from "@/components/ui/badge";
import {
  AlertCircle,
  CheckCircle,
  Clock,
  LogIn,
  LogOut,
  ShieldAlert,
} from "lucide-react";

interface SecurityEvent {
  id: string;
  eventType: string;
  ipAddress: string;
  userAgent: string;
  createdAt: string;
}

export function SecurityEvents() {
  const { toast } = useToast();
  const [events, setEvents] = useState<SecurityEvent[]>([]);
  const [isLoading, setIsLoading] = useState(true);

  useEffect(() => {
    const fetchSecurityEvents = async () => {
      try {
        const response = await fetch("/api/user/security-events", {
          method: "GET",
        });

        if (!response.ok) {
          throw new Error("Failed to fetch security events");
        }

        const data = await response.json();
        setEvents(data.events);
      } catch (error) {
        toast({
          title: "Error",
          description:
            error instanceof Error
              ? error.message
              : "Failed to fetch security events",
          variant: "destructive",
        });
      } finally {
        setIsLoading(false);
      }
    };

    fetchSecurityEvents();
  }, [toast]);

  const getEventIcon = (eventType: string) => {
    switch (eventType) {
      case "LOGIN_SUCCESS":
        return <LogIn className="h-4 w-4 text-green-500" />;
      case "LOGIN_FAILED":
        return <AlertCircle className="h-4 w-4 text-red-500" />;
      case "LOGOUT":
        return <LogOut className="h-4 w-4 text-blue-500" />;
      case "PASSWORD_CHANGED":
        return <CheckCircle className="h-4 w-4 text-green-500" />;
      case "PASSWORD_RESET_REQUESTED":
        return <Clock className="h-4 w-4 text-yellow-500" />;
      case "ACCOUNT_LOCKED":
        return <ShieldAlert className="h-4 w-4 text-red-500" />;
      default:
        return <AlertCircle className="h-4 w-4 text-gray-500" />;
    }
  };

  const getEventBadge = (eventType: string) => {
    let variant:
      | "default"
      | "secondary"
      | "destructive"
      | "outline"
      | null
      | undefined = "default";

    switch (eventType) {
      case "LOGIN_SUCCESS":
      case "PASSWORD_CHANGED":
        variant = "default";
        break;
      case "LOGOUT":
        variant = "secondary";
        break;
      case "LOGIN_FAILED":
      case "ACCOUNT_LOCKED":
        variant = "destructive";
        break;
      case "PASSWORD_RESET_REQUESTED":
        variant = "outline";
        break;
      default:
        variant = "secondary";
    }

    return (
      <Badge variant={variant} className="flex items-center gap-1">
        {getEventIcon(eventType)}
        <span>{formatEventType(eventType)}</span>
      </Badge>
    );
  };

  const formatEventType = (eventType: string) => {
    return eventType
      .replace(/_/g, " ")
      .toLowerCase()
      .replace(/\b\w/g, (char) => char.toUpperCase());
  };

  const formatDate = (dateString: string) => {
    const date = new Date(dateString);
    return new Intl.DateTimeFormat("en-US", {
      dateStyle: "medium",
      timeStyle: "short",
    }).format(date);
  };

  if (isLoading) {
    return <div>Loading security events...</div>;
  }

  if (events.length === 0) {
    return <div>No security events found.</div>;
  }

  return (
    <Table>
      <TableCaption>A list of your recent security events.</TableCaption>
      <TableHeader>
        <TableRow>
          <TableHead>Event</TableHead>
          <TableHead>IP Address</TableHead>
          <TableHead>Date & Time</TableHead>
        </TableRow>
      </TableHeader>
      <TableBody>
        {events.map((event) => (
          <TableRow key={event.id}>
            <TableCell>{getEventBadge(event.eventType)}</TableCell>
            <TableCell>{event.ipAddress}</TableCell>
            <TableCell>{formatDate(event.createdAt)}</TableCell>
          </TableRow>
        ))}
      </TableBody>
    </Table>
  );
}
