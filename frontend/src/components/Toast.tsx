import {
  createContext,
  useContext,
  useEffect,
  useState,
  type JSX,
} from "react";

interface ToastProps {
  message: string;
  type?: "info" | "success" | "warning" | "error";
  timer?: number;
  onClose: () => void;
}

interface ToastMessage {
  id: number;
  message: string;
  type?: "info" | "success" | "warning" | "error";
  timer?: number;
}

interface ToastContextType {
  addToast: (
    message: string,
    type?: "info" | "success" | "warning" | "error",
    timer?: number
  ) => void;
  addErrorToast: (error: { name: string; message: string }) => void;
}

const typeMap = {
  info: "alert-info",
  success: "alert-success",
  warning: "alert-warning",
  error: "alert-error",
};

export default function Toast({
  message,
  type = "info",
  timer = 3000,
  onClose,
}: ToastProps): JSX.Element {
  const [visible, setVisible] = useState(true);

  useEffect(() => {
    const timerId = setTimeout(() => {
      setVisible(false);
    }, timer);

    return () => {
      clearTimeout(timerId);
    };
  }, [timer]);

  useEffect(() => {
    if (!visible) {
      const timeoutId = setTimeout(() => {
        onClose();
      }, 500);

      return () => {
        clearTimeout(timeoutId);
      };
    }
  }, [visible, onClose]);

  return (
    <div
      className={`transform transition-all duration-500 ${
        visible ? "opacity-100 translate-y-0" : "opacity-0 translate-y-5"
      }`}
    >
      <div className={`alert ${typeMap[type]}`}>
        <span className="text-center text-white text-lg">{message}</span>
      </div>
    </div>
  );
}

const ToastContext = createContext<ToastContextType | undefined>(undefined);

export const useToast = (): ToastContextType => {
  const context = useContext(ToastContext);
  if (!context) {
    throw new Error("useToast must be used within a ToastProvider");
  }
  return context;
};

export const ToastProvider: React.FC<{ children: React.ReactNode }> = ({
  children,
}) => {
  const [toasts, setToasts] = useState<ToastMessage[]>([]);

  const addToast = (
    message: string,
    type: "info" | "success" | "warning" | "error" = "info",
    timer?: number
  ) => {
    const id = Date.now() + Math.random();
    setToasts((prevToasts) => [...prevToasts, { id, message, type, timer }]);
  };

  const removeToast = (id: number) => {
    setToasts((prevToasts) => prevToasts.filter((toast) => toast.id !== id));
  };

  const addErrorToast = (error: { name: string; message: string }) => {
    addToast(`${error.name}: ${error.message}`, "error");
  };

  return (
    <ToastContext.Provider value={{ addToast, addErrorToast }}>
      {children}
      <div className="toast toast-bottom toast-start flex flex-col gap-2">
        {toasts.map((toast) => (
          <Toast
            key={toast.id}
            message={toast.message}
            type={toast.type}
            timer={toast.timer}
            onClose={() => removeToast(toast.id)}
          />
        ))}
      </div>
    </ToastContext.Provider>
  );
};
