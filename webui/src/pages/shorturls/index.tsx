import ShorturlsComponent from "@/components/shorturls/Shorturls";
import { toast } from "react-toastify";
import { FormEventHandler, useCallback, useState } from "react";
import {
  Button,
  Form,
  FormFeedback,
  FormGroup,
  Input,
  Modal,
  ModalBody,
  ModalFooter,
  ModalHeader,
} from "reactstrap";
import { createUrl } from "@/api/url";
import { ShorturlData } from "@/types/ShorturlData";

const Shorturls = () => {
  const [modal, setModal] = useState(false);
  const [url, setUrl] = useState("");
  const [urlError, setUrlError] = useState("");
  const [isCreating, setIsCreating] = useState(false);
  const [shortUrls, setShortUrls] = useState<ShorturlData[]>([]);

  const stableSetShortUrls = useCallback(setShortUrls, []);
  console.log("shorturls init");

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

    const createdShorturl = await createUrl(url);
    console.log("Created short URL:", createdShorturl);
    stableSetShortUrls((prevState) => {
      // remove the same shorturl if it exists
      const newState = prevState.filter(
        (shorturl) => shorturl.id !== createdShorturl.id,
      );
      return [createdShorturl, ...newState];
    });

    setIsCreating(false);
    setModal(false);
    toast.success("Short URL created successfully", {
      position: "bottom-right",
    });
  };

  const onFormSubmit: FormEventHandler<HTMLFormElement> = (e) => {
    e.preventDefault();
    handleCreateNewShortUrl();
  };

  const onCreateClick = () => {
    handleCreateNewShortUrl();
  };

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
        <ShorturlsComponent
          shortUrls={shortUrls}
          setShortUrls={stableSetShortUrls}
        />
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
            <Form onSubmit={onFormSubmit}>
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
            <Button
              className="create"
              type="button"
              disabled={isCreating}
              onClick={onCreateClick}
            >
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
