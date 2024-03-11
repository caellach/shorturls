import ShorturlsComponent from "@/components/shorturls/Shorturls";
import { toast } from "react-toastify";
import { ChangeEvent, FormEventHandler, useState } from "react";
import {
  Button,
  ButtonProps,
  Form,
  FormFeedback,
  FormGroup,
  Input,
  Modal,
  ModalBody,
  ModalFooter,
  ModalHeader,
} from "reactstrap";
import { ShorturlData } from "@/types/ShorturlData";

const Shorturls = () => {
  const [modal, setModal] = useState(false);
  const [url, setUrl] = useState("");
  const [urlError, setUrlError] = useState("");
  const [isCreating, setIsCreating] = useState(false);
  const [shortUrls, setShortUrls] = useState<ShorturlData[]>([]);

  const toggleModalCreateNewShortUrl = () => {
    setModal(!modal);
    setTimeout(() => {
      setUrl("");
    }, 500);
  };

  const handleCreateNewShortUrl = async () => {
    if (!url) {
      setUrlError("URL is required");
      return;
    }

    if (!/^https?:\/\/[^/\s]+\/?\S*$/.test(url)) {
      setUrlError("URL is invalid");
      return;
    }

    console.log("URL:", url);
    // Handle the creation of the new short URL here
    setIsCreating(true);

    // use a promise to delay the response
    const promise = new Promise<void>((resolve) => {
      setTimeout(() => {
        resolve();
      }, 1000);
    });
    await Promise.all([promise]);

    setShortUrls([
      ...shortUrls,
      {
        id: "123",
        url,
        shortUrl: "http://localhost:3000/123",
        useCount: 0,
        lastUsed: new Date(0),
      },
    ]);

    setIsCreating(false);
    setModal(false);
    toast.success("Short URL created successfully");
  };

  const handleFormSubmit = (e: ButtonProps) => {};

  return (
    <>
      <div className="shorturls-content-wrapper">
        <div className="shorturls-title-wrapper">
          <h1>Short URLs</h1>
          <div>
            <button
              type="button"
              className="create"
              onClick={toggleModalCreateNewShortUrl}
            >
              Create +
            </button>
          </div>
        </div>
        <ShorturlsComponent shortUrls={shortUrls} setShortUrls={setShortUrls} />
        <Modal
          centered
          isOpen={modal}
          toggle={toggleModalCreateNewShortUrl}
          backdrop={isCreating ? "static" : true}
        >
          <ModalHeader toggle={toggleModalCreateNewShortUrl}>
            <b>Create New Short URL</b>
          </ModalHeader>
          <ModalBody>
            <Form
              onSubmit={(e) => {
                e.preventDefault();
                handleCreateNewShortUrl();
              }}
            >
              <FormGroup>
                <Input
                  type="url"
                  value={url}
                  onChange={(e) => {
                    setUrl(e.target.value);
                    setUrlError("");
                  }}
                  placeholder="Enter URL"
                  invalid={!!urlError}
                />
                <FormFeedback>{urlError}</FormFeedback>
              </FormGroup>
            </Form>
          </ModalBody>
          <ModalFooter>
            <Button className="create" type="submit" disabled={isCreating}>
              Create
            </Button>
            <Button
              color="secondary"
              type="reset"
              onClick={toggleModalCreateNewShortUrl}
              disabled={isCreating}
            >
              Cancel
            </Button>
          </ModalFooter>
        </Modal>
      </div>
    </>
  );
};

export default Shorturls;
